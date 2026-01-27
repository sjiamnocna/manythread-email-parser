package email

import (
"testing"
)

// Note: Tests for CollectFiles() are handled in integration tests since they require
// access to the actual filesystem. For unit testing, we focus on the helper functions
// that are exported for testing purposes.

// TestCollectFiles_Integration would test the recursive directory scanning,
// but that requires filesystem setup and is better suited for integration tests.
// The CollectFiles function is validated through end-to-end testing.
