package main

import "fmt"

func main() {
	deployer := &Deployer{
		cfg:           cfg,
		Destination:   args.out,
		Version:       args.version,
		OutputSource:  args.src,
		OutputSiteOmw: args.siteOmw,
	}

	err := deployer.Deploy()
	if err != nil {
		fmt.Println("错误:", err)
	} else {
		fmt.Println("成功,输出目录:", deployer.Destination)
	}
}
