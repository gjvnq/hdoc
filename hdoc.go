package main

import (
	"encoding/xml"
	"fmt"
	"github.com/beevik/etree"
	"github.com/spf13/cobra"
	"os"
	"regexp"
	"strings"
)

var reWhitespace = regexp.MustCompile(`\s+`)
var reInvalidId = regexp.MustCompile(`[^\p{L}\p{N}_-]+`)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hdoc [source]",
	Args:  cobra.ExactArgs(1),
	Short: "A simple HTML shorthand processor",
	// 	Long: `A longer description that spans multiple lines and likely contains
	// examples and usage of using your application. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	Run: realMain,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func init() {
	// cobra.OnInitialize(initConfig)
	rootCmd.Flags().StringP("output", "o", "", "Output file")
	rootCmd.Flags().StringP("style", "", "", "Set document style options")
}

func realMain(cmd *cobra.Command, args []string) {
	var outFile *os.File
	var err error

	// Open output
	outfn := cmd.Flag("output").Value.String()
	if outfn != "" {
		outFile, err = os.OpenFile(outfn, os.O_RDWR|os.O_CREATE, 0644)
		panicIfErr(err)
	} else {
		outFile = os.Stdout
	}
	fmt.Println(outFile)

	// Open style
	stylefn := cmd.Flag("style").Value.String()
	if stylefn == "" {
		stylefn = "data/styles/default.html"
	}
	style := etree.NewDocument()
	style.ReadSettings.Entity = xml.HTMLEntity
	err = style.ReadFromFile(stylefn)
	panicIfErr(err)
	fmt.Println(style)

	// Open source
	src := etree.NewDocument()
	style.ReadSettings.Entity = xml.HTMLEntity
	err = src.ReadFromFile(args[0])
	panicIfErr(err)
	fmt.Println(src)

	fmt.Println(getInnerTxt(src.Root()))

	processAbbrAndDfn(src)

	src.WriteTo(outFile)
}

func getInnerTxt(cur *etree.Element) string {
	buf := new(strings.Builder)
	getInnerTxt_inner(cur, buf)
	clean := reWhitespace.ReplaceAllString(buf.String(), " ")
	return strings.TrimSpace(clean)
}

func getInnerTxt_inner(cur *etree.Element, buf *strings.Builder) {
	for _, raw_child := range cur.Child {
		switch child := raw_child.(type) {
		case *etree.CharData:
			buf.WriteString(child.Data)
		case *etree.Element:
			getInnerTxt_inner(child, buf)
		}
	}
}

func fixWAttr(elem *etree.Element) {
	if elem.SelectAttr("w") != nil {
		return
	}
	txt := getInnerTxt(elem)
	elem.CreateAttr("w", txt)
}

func fixIdAttr(elem *etree.Element) {
	id := elem.SelectAttrValue("id", "")
	if id != "" {
		return
	}
	w := elem.SelectAttrValue("w", "")
	w = reInvalidId.ReplaceAllString(w, "")
	elem.CreateAttr("id", elem.Tag+"-"+w)
}

func classifyAbbrDfn(elem *etree.Element) string {
	if elem.Space == "" && (elem.Tag == "dfn" || elem.Tag == "abbr") {
		if len(elem.Child) == 0 {
			return "ref"
		}
		return "def"
	}
	return "neither"
}

func processAbbrAndDfn(src *etree.Document) {
	m := make(map[string]*etree.Element)

	processAbbrAndDfn_inner(src.Root(), m)
}

func addClass(elem *etree.Element, class string) {
	val := elem.SelectAttrValue("class", "")
	list := strings.Split(val, " ")
	// Hack to deal with empty strings
	if len(list) > 0 && len(list[0]) == 0 {
		list = list[1:]
	}
	for _, item := range list {
		if item == class {
			return
		}
	}
	list = append(list, class)
	val = strings.Join(list, " ")
	elem.CreateAttr("class", val)
}

func processAbbrAndDfn_inner(cur *etree.Element, m map[string]*etree.Element) {
	for _, raw_child := range cur.Child {
		switch child := raw_child.(type) {
		case *etree.Element:
			kind := classifyAbbrDfn(child)
			if kind == "def" {
				fixWAttr(child)
				if child.Tag == "dfn" {
					fixIdAttr(child)
				}
				w_attr := child.Tag + "-" + child.SelectAttrValue("w", "")
				m[w_attr] = child
				child.RemoveAttr("w")
			}
			if kind == "ref" {
				fixWAttr(child)
				w_attr := child.Tag + "-" + child.SelectAttrValue("w", "")
				def := m[w_attr].Copy()
				if child.Tag == "dfn" {
					child.Tag = "a"
					child.CreateAttr("href", "#"+def.SelectAttrValue("id", ""))
					child.Child = def.Child
					addClass(child, "dfn")
				}
				if child.Tag == "abbr" {
					title := def.SelectAttrValue("title", "")
					child.Child = def.Child
					child.CreateAttr("title", title)
				}
				child.RemoveAttr("w")
			}

			processAbbrAndDfn_inner(child, m)
		}
	}
}
