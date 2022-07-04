package inventory

import "github.com/sergi/go-diff/diffmatchpatch"

func makePatch(text1, text2 string) string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(text1, text2, false)
	patchs := dmp.PatchMake(text1, text2, diffs)
	return dmp.PatchToText(patchs)
}

func patch(text1, patch string) (string, error) {
	dmp := diffmatchpatch.New()
	patches, err := dmp.PatchFromText(patch)
	if err != nil {
		return "", err
	}
	text2, _ := dmp.PatchApply(patches, text1)
	return text2, nil

}

func patches(text1 string, patcheStrs ...string) (string, error) {
	dmp := diffmatchpatch.New()
	var patches []diffmatchpatch.Patch
	for _, patch := range patcheStrs {
		p, err := dmp.PatchFromText(patch)
		if err != nil {
			return "", err
		}
		patches = append(patches, p...)
	}
	text2, _ := dmp.PatchApply(patches, text1)
	return text2, nil
}
