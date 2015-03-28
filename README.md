# Imstor
A Golang image storage engine. Used to create and store different sizes/thumbnails of user uploaded images.

[![Build Status](https://travis-ci.org/deiwin/imstor.svg?branch=master)](https://travis-ci.org/deiwin/imstor)
[![Coverage](http://gocover.io/_badge/github.com/deiwin/imstor?0)](http://gocover.io/github.com/deiwin/imstor)
[![GoDoc](https://godoc.org/github.com/deiwin/imstor?status.svg)](https://godoc.org/github.com/deiwin/imstor)

## Description

**Imstor** enables you to create copies (or thumbnails) of your images and stores
them along with the original image on your filesystem. The image and its
copies are stored in a file structure based on the (zero-prefixed, decimal)
CRC 64 checksum of the original image. The last 2 characters of the checksum
are used as the lvl 1 directory name.

**Imstor** supports any image format you can decode to go's own image.Image
and then back to your preferred format. The decoder for any given image is
chosen by the image's mimetype.

### Example folder structure
Folder name and contents, given the checksum `08446744073709551615` and
sizes named "*small*" and "*large*":
```
/configured/root/path/15/08446744073709551615/original.jpg
/configured/root/path/15/08446744073709551615/small.jpg
/configured/root/path/15/08446744073709551615/large.jpg
```

## Usage
See tests for usage examples.
