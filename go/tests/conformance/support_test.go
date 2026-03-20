package conformance

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

type conformanceExpected struct {
	Result json.RawMessage `json:"result"`
	Error  *struct {
		Message string `json:"message"`
	} `json:"error"`
}

type languageException struct {
	Reason string          `json:"reason"`
	Result json.RawMessage `json:"result"`
	Error  *struct {
		Message string `json:"message"`
	} `json:"error"`
	Skip bool `json:"skip"`
}

type resolvedExpectation struct {
	Result json.RawMessage
	Error  string
	Skip   string
}

type conformanceCase struct {
	ID       string
	Input    json.RawMessage
	Expected resolvedExpectation
}

func loadConformanceCases(t *testing.T, suite string) []conformanceCase {
	t.Helper()

	inputRoot := filepath.Join(repoRoot(t), "test", "input", suite)
	expectedRoot := filepath.Join(repoRoot(t), "test", "expected", suite)

	var cases []conformanceCase
	for _, inputPath := range walkJSONFiles(t, inputRoot) {
		relativePath, err := filepath.Rel(inputRoot, inputPath)
		if err != nil {
			t.Fatalf("rel %q: %v", inputPath, err)
		}

		expectedPath := filepath.Join(expectedRoot, relativePath)
		if _, err := os.Stat(expectedPath); err != nil {
			t.Fatalf("missing expected file for %s/%s: %v", suite, relativePath, err)
		}

		testID := trimJSONSuffix(filepath.ToSlash(relativePath))
		inputPayload := mustReadJSONRaw(t, inputPath)
		expectedPayload := mustReadExpected(t, expectedPath)
		cases = append(cases, conformanceCase{
			ID:       testID,
			Input:    inputPayload,
			Expected: resolveExpectation(t, suite, testID, expectedPayload),
		})
	}

	return cases
}

func repoRoot(t *testing.T) string {
	t.Helper()

	root, err := filepath.Abs(filepath.Join("..", "..", ".."))
	if err != nil {
		t.Fatalf("resolve repo root: %v", err)
	}
	return root
}

func walkJSONFiles(t *testing.T, root string) []string {
	t.Helper()

	var paths []string
	if err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".json" {
			paths = append(paths, path)
		}
		return nil
	}); err != nil {
		t.Fatalf("walk %q: %v", root, err)
	}

	sort.Strings(paths)
	return paths
}

func mustReadJSONRaw(t *testing.T, path string) json.RawMessage {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %q: %v", path, err)
	}

	var raw json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("decode %q: %v", path, err)
	}
	return raw
}

func mustReadExpected(t *testing.T, path string) conformanceExpected {
	t.Helper()

	var expected conformanceExpected
	decodeJSONFile(t, path, &expected)
	return expected
}

func resolveExpectation(t *testing.T, suite, caseID string, expected conformanceExpected) resolvedExpectation {
	t.Helper()

	exceptionPath := filepath.Join(repoRoot(t), "test", "language-exceptions", "go", suite, caseID+".json")
	if _, err := os.Stat(exceptionPath); err == nil {
		var exception languageException
		decodeJSONFile(t, exceptionPath, &exception)

		if exception.Reason == "" {
			t.Fatalf("language exception %q is missing a reason", suite+"/"+caseID)
		}
		if exception.Skip {
			if len(exception.Result) != 0 || exception.Error != nil {
				t.Fatalf("language exception %q cannot define skip with result or error", suite+"/"+caseID)
			}
			return resolvedExpectation{Skip: exception.Reason}
		}

		hasResult := len(exception.Result) != 0
		hasError := exception.Error != nil
		if hasResult == hasError {
			t.Fatalf("language exception %q must define exactly one of result or error", suite+"/"+caseID)
		}

		resolved := resolvedExpectation{Result: exception.Result}
		if exception.Error != nil {
			resolved.Error = exception.Error.Message
		}
		return resolved
	}

	hasResult := len(expected.Result) != 0
	hasError := expected.Error != nil
	if hasResult == hasError {
		t.Fatalf("conformance case %q must define exactly one of result or error", suite+"/"+caseID)
	}

	resolved := resolvedExpectation{Result: expected.Result}
	if expected.Error != nil {
		resolved.Error = expected.Error.Message
	}
	return resolved
}

func decodeJSONFile(t *testing.T, path string, destination any) {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %q: %v", path, err)
	}
	if err := json.Unmarshal(data, destination); err != nil {
		t.Fatalf("decode %q: %v", path, err)
	}
}

func trimJSONSuffix(path string) string {
	const suffix = ".json"
	return path[:len(path)-len(suffix)]
}

func runConformanceSuite[TInput any, TResult any](
	t *testing.T,
	suite string,
	run func(TInput) (TResult, error),
) {
	t.Helper()

	for _, testCase := range loadConformanceCases(t, suite) {
		t.Run(testCase.ID, func(t *testing.T) {
			if testCase.Expected.Skip != "" {
				t.Skip(testCase.Expected.Skip)
			}

			input := decodeRawJSON[TInput](t, testCase.Input)
			result, err := run(input)

			if testCase.Expected.Error != "" {
				if err == nil {
					t.Fatalf("expected error %q, got nil", testCase.Expected.Error)
				}
				if err.Error() != testCase.Expected.Error {
					t.Fatalf("expected error %q, got %q", testCase.Expected.Error, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			expected := decodeRawJSON[TResult](t, testCase.Expected.Result)
			if !reflect.DeepEqual(result, expected) {
				t.Fatalf("unexpected result:\nexpected: %#v\ngot: %#v", expected, result)
			}
		})
	}
}

func decodeRawJSON[T any](t *testing.T, raw json.RawMessage) T {
	t.Helper()

	var value T
	if err := json.Unmarshal(raw, &value); err != nil {
		t.Fatalf("decode raw json: %v", err)
	}
	return value
}
