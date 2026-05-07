package base

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/jessevdk/go-flags"
	"golang.org/x/image/bmp"
)

// InputCommand is the base command for commands that take an input file.
type InputCommand struct {
	// Input is the name of the input file.
	Input flags.Filename `short:"i" long:"input" description:"The name of the input file or - for STDIN" optional:"true" default:"-"`
}

// ReadInput reads the input image from the input stream.
func (cmd *InputCommand) ReadInput() (image.Image, error) {
	// open the input stream
	var input io.Reader

	if cmd.Input == "-" {
		// read from standard input
		slog.Debug("reading input from STDIN")
		input = os.Stdin
	} else {
		// open the underlay image file
		slog.Debug("reading input from file", "name", cmd.Input)
		var err error
		if input, err = os.Open(string(cmd.Input)); err != nil {
			slog.Error("error opening input file", "name", cmd.Input, "error", err)
			return nil, err
		}
		defer input.(io.ReadCloser).Close()
	}

	// decode the base image
	base, _, err := image.Decode(input)
	if err != nil {
		slog.Error("error decoding input data for base image", "name", cmd.Input, "error", err)
		return nil, err
	}

	return base, nil
}

// OutputCommand is the base command for commands that produce an output file.
type OutputCommand struct {
	// Output is the name of the output file.
	Output flags.Filename `short:"o" long:"output" description:"The name of the output file or - for STDOUT" optional:"true" default:"-"`
	// Format is the output format, if an output filename is not specified; it is used for chaining.
	Format string `short:"x" long:"format" description:"Format of the output image" optional:"true" choice:"jpeg" choice:"jpg" choice:"png" choice:"gif" choice:"bmp" default:"png"`
	// DPI is the image resolution in Dots Per Inch.
	DPI float64 `short:"d" long:"dpi" description:"The image resolution in DPI - Dots Per Inch" optional:"true" default:"72"`
}

// OutputStream returns an io.Writer for the output file.
func (cmd *OutputCommand) WriteOutput(img image.Image) error {
	// open the output stream
	var (
		output io.Writer
		err    error
	)

	if cmd.Output == "-" {
		// writing image to standard output
		slog.Debug("writing image to STDOUT", "format", cmd.Format)
		output = os.Stdout
	} else {
		// check the output format
		switch strings.ToLower(filepath.Ext(string(cmd.Output))) {
		case ".jpg", ".jpeg":
			cmd.Format = "jpg"
		case ".png":
			cmd.Format = "png"
		case ".gif":
			cmd.Format = "gif"
		case ".bmp":
			cmd.Format = "bmp"
		default:
			err = fmt.Errorf("unsupported output file type: %s", filepath.Ext(string(cmd.Output)))
			slog.Error("unsupported output image type", "name", cmd.Output)
			return err
		}
		slog.Debug("writing output to file", "name", cmd.Output, "format", cmd.Format)

		// open the output file
		if output, err = os.Create(string(cmd.Output)); err != nil {
			slog.Error("error opening output file", "name", cmd.Output, "error", err)
			return err
		}
		defer output.(io.WriteCloser).Close()
	}

	// encode the output image
	slog.Debug("encoding output image", "name", cmd.Output, "format", cmd.Format)
	switch cmd.Format {
	case "jpg", "jpeg":
		slog.Debug("encoding output file as JPEG", "name", cmd.Output)
		if err = jpeg.Encode(output, img, nil); err != nil {
			slog.Error("error encoding output file", "name", cmd.Output, "error", err, "format", cmd.Format)
			return err
		}
	case "png":
		slog.Debug("encoding output file as PNG", "name", cmd.Output)
		if err = png.Encode(output, img); err != nil {
			slog.Error("error encoding output file", "name", cmd.Output, "error", err, "format", cmd.Format)
			return err
		}
	case "gif":
		slog.Debug("encoding output file as GIF", "name", cmd.Output)
		if err = gif.Encode(output, img, nil); err != nil {
			slog.Error("error encoding output file", "name", cmd.Output, "error", err, "format", cmd.Format)
			return err
		}
	case "bmp":
		slog.Debug("encoding output file as BMP", "name", cmd.Output)
		if err = bmp.Encode(output, img); err != nil {
			slog.Error("error encoding output file", "name", cmd.Output, "error", err, "format", cmd.Format)
			return err
		}
	default:
		fmt.Fprintf(os.Stderr, "Unsupported output format: %s\n", cmd.Format)
		slog.Error("unsupported output format", "name", cmd.Output, "format", cmd.Format)
		return fmt.Errorf("unsupported output format: %s", cmd.Format)
	}

	slog.Debug("output image written")
	return nil
}
