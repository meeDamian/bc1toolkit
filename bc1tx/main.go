package main

var Opts struct {
	Blab string `short:"b" long:"blab" description:"Blab blab blab"`
}

func init() {
	//common.Define("bla", &Opts)
	//common.Parse()
}

func main() {

	//			func() string {
	//	return fmt.Sprintf("Usage of %s:\n", os.Args[0])
	//}

	//stat, _ := os.Stdin.Stat()
	//if (stat.Mode() & os.ModeCharDevice) == 0 {
	//	fmt.Println("data is being piped to stdin")
	//} else {
	//	fmt.Println("stdin is from a terminal")
	//}
	//
	//buf := &bytes.Buffer{}
	//n, err := io.Copy(buf, os.Stdin)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//if n <= 1 { // buffer always contains '\n'
	//	log.Fatalln("no input provided")
	//}
	//
	//if len(os.Args) != 3 {
	//	log.Fatalln("usage: echo \"hello world\" | change hello bye")
	//}
	//
	//oldWord := os.Args[1]
	//newWord := os.Args[2]
	//
	//r := bytes.Replace(buf.Bytes(), []byte(oldWord), []byte(newWord), -1)
	//
	//fmt.Println(string(r))
}
