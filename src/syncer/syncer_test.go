package syncer

import "testing"

// CreateAProjectFromCheckout tests

func TestCanOpenAProjectCheckout(t *testing.T) {
	projectName := "Foobar"
	project := CreateProjectFromCheckout(projectName)

	expectedName := "Foobar"
	if project.Name != expectedName {
		t.Logf("%v is not %v", project.Name, expectedName)
		t.Fail()
	}
}
