package manifest

import "testing"

func TestManifestCanBeParseable(t *testing.T) {
	manifest, _ := ParseManifest("../../testdata/fixtures/example_project/manifest.yml")

	var (
		expectedName     = "ExampleProject"
		expectedOwner    = "developer-team@example.com"
		expectedLanguage = "go"
	)

	if manifest.Name != expectedName {
		t.Logf("%v is not %v", manifest.Name, expectedName)
		t.Fail()
	}

	if manifest.Owner != expectedOwner {
		t.Logf("%v is not %v", manifest.Owner, expectedOwner)
		t.Fail()
	}

	if manifest.Language != expectedLanguage {
		t.Logf("%v is not %v", manifest.Owner, expectedLanguage)
		t.Fail()
	}
}
