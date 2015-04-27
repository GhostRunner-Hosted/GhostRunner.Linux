package utils

func StringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func SliceToList(sliceList []string) (string) {
    var list string

    for i := 0; i < len(sliceList); i++ {
    	if (i != 0) {
    		list += ","
    	}

    	list += sliceList[i]
    }

    return list
}