
This is a fork of pdfcpu/pdfcpu, a lightweight version without unnecessary files that you don't use when working with its API

Project link: pdfcpu.io/

Many tests are deleted because they depend on test pdfs and fonts, dependencies are replaced:

    github.com/pdfcpu/pdfcpu/ -> path/to/pdfcpu

In the main.go you can test merging of two pdf files:

    make build
    bin/example_merge 
The result file will be created in the ` cmd/example_merge` folder
