package huffman_test

import (
    "huffman"
    "testing"
    "fmt"
)

func TestBuildHuffman(t *testing.T) {

    runes := []byte{'A', 'B', 'C', 'D'}

    f := make(map[byte]int)
    f[runes[0]] = 3
    f[runes[1]] = 2
    f[runes[2]] = 1
    f[runes[3]] = 1
    
    hf, err := huffman.BuildHuffman(f)
    if err != nil {
        t.Error("err was not nil", err)
    }

    expected := &huffman.Huffman{
        nil,
        7,
        &huffman.Huffman{
            &runes[0],
            3,
            nil,
            nil,
        },
        &huffman.Huffman{
            nil,
            4,
            &huffman.Huffman{
                &runes[1],
                2,
                nil,
                nil,
            },
            &huffman.Huffman{
                nil,
                2,
                &huffman.Huffman{
                    &runes[2],
                    1,
                    nil,
                    nil,
                },
                &huffman.Huffman{
                    &runes[3],
                    1,
                    nil,
                    nil,
                },
            },
        },
    }

    check := huffman.CompareHuffman(hf, expected)
    if !check {
        huffman.Print(hf, "actual")
        huffman.Print(expected, "expected")

        t.Error("Trees are not equal")
        t.FailNow()
    }

}

func TestHuffmanBuildTable(t *testing.T) {
    runes := []byte{'A', 'B', 'C', 'D'}

    f := make(map[byte]int)
    f[runes[0]] = 3
    f[runes[1]] = 2
    f[runes[2]] = 1
    f[runes[3]] = 1
    
    hf, err := huffman.BuildHuffman(f)
    if err != nil {
        t.Error("err was not nil", err)
        t.FailNow()
    }

    enc, err := huffman.GetTable(hf)
    if err != nil {
        t.Error("err was nil", err)
        t.FailNow()
    }
    expectedEnc := make(map[byte]uint32)
    //expectedEnc['A'] = "0" // 0
    //expectedEnc['B'] = "10" // 2
    //expectedEnc['C'] = "110" // 6
    //expectedEnc['D'] = "111" // 7
    expectedEnc['A'] = 0
    expectedEnc['B'] = 2
    expectedEnc['C'] = 6
    expectedEnc['D'] = 7

    if len(enc) != len(expectedEnc) {
        t.Error("len(enc) != len(expectedEnc)")
        t.FailNow()
    }

    for k, v := range enc {
        if (v != expectedEnc[k]) {
            t.Errorf("Wanted %v, got %v", expectedEnc[k], v)
            fmt.Printf("Actual: %+v\n", enc)
            fmt.Printf("Expected: %+v\n", expectedEnc)
            t.FailNow()
        }
    }
}

func TestHuffmanEncode(t *testing.T) {
    runes := []byte{'A', 'B', 'C', 'D'}

    f := make(map[byte]int)
    f[runes[0]] = 3
    f[runes[1]] = 2
    f[runes[2]] = 1
    f[runes[3]] = 1
    
    hf, _ := huffman.BuildHuffman(f)
    codes, err := huffman.GetTable(hf)
    if err != nil {
        t.Error("err was nil", err)
        t.FailNow()
    }

    test := []byte{'A', 'B', 'B', 'C', 'C', 'D', 'A', 'B'}
    actual, err := huffman.Encode(codes, test)
    if err != nil {
        t.Error("err was nil", err)
        t.FailNow()
    }
    expected := &huffman.EncodedData{
        codes,
        []uint8{86, 221, 0},
        7,
    }

    actualStr := fmt.Sprintf("%+v", actual)
    expectedStr := fmt.Sprintf("%+v", expected)
    if actualStr != expectedStr {
        t.Errorf("wanted:\n\n%s\n\ngot:\n\n%s\n", expectedStr, actualStr)
        t.FailNow()
    }
}
