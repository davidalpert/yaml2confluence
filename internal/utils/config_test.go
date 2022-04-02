package utils

type MockFileSystem struct {
	wd string
}

func (mfs MockFileSystem) Getwd() (string, error) { return mfs.wd, nil }

// func TestResolveAbsolutePath(t *testing.T) {

// 	fs = MockFileSystem{"/home/joe"}

// 	path := ResolveAbsolutePath("../test")
// 	test, _ := filepath.Abs("../../test")
// 	t.Fatal("TEST " + test)
// 	fmt.Println(path)
// }
