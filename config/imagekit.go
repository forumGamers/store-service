package config

import (
	"errors"
	"os"

	"context"

	"github.com/codedius/imagekit-go"
	"github.com/joho/godotenv"
)

func ImagekitConnection() *imagekit.Client{
	if err := godotenv.Load() ; err != nil {
		panic(err.Error())
	}

	IMAGEKIT_PRIVATE_KEY := os.Getenv("IMAGEKIT_PRIVATE_KEY")
	IMAGEKIT_PUBLIC_KEY := os.Getenv("IMAGEKIT_PUBLIC_KEY")

	opts := imagekit.Options{
		PublicKey: IMAGEKIT_PUBLIC_KEY,
		PrivateKey: IMAGEKIT_PRIVATE_KEY,
	}

	if ik,err := imagekit.NewClient(&opts) ; err != nil {
		panic(err.Error())
	}else {
		return ik
	}
}

func UploadImage(file []byte,fileName string) (string, string , error){

	ur := imagekit.UploadRequest{
		File: file,
		FileName: fileName,
	}

	ctx := context.Background()

	if upr,err := ImagekitConnection().Upload.ServerUpload(ctx,&ur) ; err != nil {
		return "",
			   "",
			   errors.New(err.Error())
	}else {
		return upr.URL,
			   upr.FileID,
			   nil
	}
}