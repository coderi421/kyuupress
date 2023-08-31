package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func getGhostscriptPath() string {
	gsNames := []string{"gs", "gswin32", "gswin64"}
	for _, name := range gsNames {
		if exec.Command("sh", "-c", "command -v "+name).Run() == nil {
			return name
		}
	}
	return ""
}

func compressPDF(inputFilePath, outputFilePath string, power int) error {
	gsPath := getGhostscriptPath()
	if gsPath == "" {
		return fmt.Errorf("未找到 Ghostscript")
	}

	quality := map[int]string{
		0: "/default",
		1: "/prepress",
		2: "/printer",
		3: "/ebook",
		4: "/screen",
	}

	cmd := exec.Command(gsPath,
		"-sDEVICE=pdfwrite",
		"-dCompatibilityLevel=1.4",
		"-dPDFSETTINGS="+quality[power],
		"-dNOPAUSE", "-dQUIET", "-dBATCH",
		"-sOutputFile="+outputFilePath,
		inputFilePath)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("压缩 PDF 文件时出现错误: %v\n%s", err, output)
	}

	return nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("请输入输入的 PDF 文件路径: ")
	inputPath, _ := reader.ReadString('\n')
	inputPath = strings.TrimSpace(inputPath)

	fmt.Print("请输入输出的压缩后 PDF 文件路径: ")
	outputPath, _ := reader.ReadString('\n')
	outputPath = strings.TrimSpace(outputPath)

	fmt.Print("请输入压缩比例（0-4，0 表示默认压缩）: ")
	ratioStr, _ := reader.ReadString('\n')
	ratioStr = strings.TrimSpace(ratioStr)
	compressionRatio, err := strconv.Atoi(ratioStr)
	if err != nil || compressionRatio < 0 || compressionRatio > 4 {
		fmt.Println("无效的压缩比例")
		return
	}

	// 创建输出文件夹
	outputFolder := filepath.Dir(outputPath)
	if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
		err := os.Mkdir(outputFolder, 0755)
		if err != nil {
			fmt.Println("无法创建输出文件夹:", err)
			return
		}
	}

	err = compressPDF(inputPath, outputPath, compressionRatio)
	if err != nil {
		fmt.Println("压缩 PDF 文件时出现错误:", err)
		return
	}

	fmt.Println("PDF 文件压缩成功！")
}
