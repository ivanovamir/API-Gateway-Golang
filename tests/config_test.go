package tests

import (
	"gateway/pkg/config"
	"testing"
)

func TestConfig(t *testing.T) {
	testTable := []struct {
		name            string
		configPath      string
		expectError     bool
		expectErrorBody string
	}{
		{
			name:        "test_1",
			configPath:  "configs/test1.json",
			expectError: false,
		},
		{
			name:            "test_2",
			configPath:      "configs/test2.json",
			expectError:     true,
			expectErrorBody: "proxy url is empty",
		},
		{
			name:            "test_3",
			configPath:      "configs/test3.json",
			expectError:     true,
			expectErrorBody: "proxy_method is empty",
		},
		{
			name:            "test_4",
			configPath:      "configs/test4.json",
			expectError:     true,
			expectErrorBody: "proxy_method must be uppercase",
		},
		{
			name:            "test_5",
			configPath:      "configs/test5.json",
			expectError:     true,
			expectErrorBody: "expected_proxy_status_codes is empty",
		},
		{
			name:       "test_6",
			configPath: "configs/test6.json",
		},
		{
			name:            "test_7",
			configPath:      "configs/test7.json",
			expectError:     true,
			expectErrorBody: "method is empty",
		},
		{
			name:            "test_8",
			configPath:      "configs/test8.json",
			expectError:     true,
			expectErrorBody: "path is empty",
		},
		{
			name:            "test_9",
			configPath:      "configs/test9.json",
			expectError:     true,
			expectErrorBody: "url is empty",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			_, err := config.NewConfig(config.WithPath(testCase.configPath)).ParseConfig()

			if testCase.expectError {
				if testCase.expectErrorBody != err.Error() {
					t.Errorf("error occured parsing config: %v", err)
				}
			}
		})
	}
}
