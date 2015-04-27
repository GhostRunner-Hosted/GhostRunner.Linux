package utils

import (
    "strings"
    "strconv"
    "crypto/aes"
    "crypto/cipher"
    "math/rand"
    "encoding/hex"
    "errors"
    "ghostrunner/logging"
)

func Encrypt(key, text string) (string) {
    logging.Debug("encryption", "Decrypt", "Starting to decrypt task log");

    ciphertext := []byte(text)
    ciphertext = pad(ciphertext)

    if len(ciphertext)%aes.BlockSize != 0 {
        logging.Error("encryption", "Encrypt", "Block sizes are incorrect", errors.New("Cipher text is " + strconv.Itoa(len(ciphertext)) + " long but should be " + strconv.Itoa(aes.BlockSize)))
    } else {
        block, err := aes.NewCipher([]byte(key))
        
        if err != nil {
            logging.Error("encryption", "Encrypt", "Error creating cipher", err)
        } else {
            iv := []byte(generateIV(16))
            
            mode := cipher.NewCBCEncrypter(block, iv)

            mode.CryptBlocks(ciphertext, ciphertext)

            return  hex.EncodeToString(ciphertext) + "$" + hex.EncodeToString(iv)
        }
    }
   
    return ""
}

func Decrypt(key, text string) ([]byte) {
    logging.Debug("encryption", "Decrypt", "Starting to decrypt task script");

    var ciphertext []byte

    textParts := strings.Split(text, "$")

    logging.Debug("encryption", "Decrypt", "Checking the encrypted script contains all parts");

    if (len(textParts) == 3) {
        logging.Debug("encryption", "Decrypt", "All parts are located");

        ciphertext, _ = hex.DecodeString(textParts[0])

        block, _ := aes.NewCipher([]byte(key))

        if len(ciphertext)%aes.BlockSize != 0 {
            logging.Error("encryption", "Decrypt", "Block sizes are incorrect", errors.New("Cipher text is " + strconv.Itoa(len(ciphertext)) + " long but should be " + strconv.Itoa(aes.BlockSize)))
        } else {
            logging.Debug("encryption", "Decrypt", "Text looks good, decrypting");

            mode := cipher.NewCBCDecrypter(block, []byte(textParts[1]))

            mode.CryptBlocks(ciphertext, ciphertext)

            ciphertext = unpad(ciphertext)
        }
    } else {
        logging.Error("encryption", "Decrypt", "The encrypted script does not contain all parts", errors.New("Encrypted script only contains " + strconv.Itoa(len(textParts)) + " parts"))
    }
    
    return ciphertext
}

func generateIV(n int) string {
    var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func pad(in []byte) []byte {
    padding := 16 - (len(in) % 16)
    
    if padding == 0 {
        padding = 16
    }
    
    for i := 0; i < padding; i++ {
        in = append(in, byte(padding))
    }

    return in
}

func unpad(in []byte) []byte {
    if len(in) == 0 {
        return nil
    }

    padding := in[len(in)-1]
    
    if int(padding) > len(in) || padding > aes.BlockSize {
        return nil
    } else if padding == 0 {
        return nil
    }

    for i := len(in) - 1; i > len(in)-int(padding)-1; i-- {
        if in[i] != padding {
        return nil
        }
    }

    return in[:len(in)-int(padding)]
}