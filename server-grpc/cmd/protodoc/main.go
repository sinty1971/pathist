package main

import (
	"fmt"
	"sort"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	pluginpb "google.golang.org/protobuf/types/pluginpb"
)

type docConfig struct {
	title    string
	filename string
}

func main() {
	cfg := docConfig{
		title:    "API Documentation",
		filename: "apis.md",
	}

	opts := protogen.Options{
		ParamFunc: func(name, value string) error {
			switch name {
			case "title":
				if value != "" {
					cfg.title = value
				}
			case "filename":
				if value != "" {
					cfg.filename = value
				}
			default:
				return fmt.Errorf("unknown parameter %q", name)
			}
			return nil
		},
	}

	opts.Run(func(plugin *protogen.Plugin) error {
		plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_SUPPORTS_EDITIONS)

		if cfg.filename == "" {
			cfg.filename = "apis.md"
		}

		files := collectFiles(plugin.Files)
		generated := plugin.NewGeneratedFile(cfg.filename, "")

		fmt.Fprintf(generated, "# %s\n\n", cfg.title)
		if len(files) == 0 {
			fmt.Fprintln(generated, "_生成対象の proto ファイルがありません。_")
			return nil
		}

		for _, file := range files {
			writeFileDoc(generated, file)
		}

		return nil
	})
}

func collectFiles(files []*protogen.File) []*protogen.File {
	var targets []*protogen.File
	for _, file := range files {
		if file.Generate {
			targets = append(targets, file)
		}
	}
	sort.Slice(targets, func(i, j int) bool {
		return targets[i].Desc.Path() < targets[j].Desc.Path()
	})
	return targets
}

func writeFileDoc(g *protogen.GeneratedFile, file *protogen.File) {
	fmt.Fprintf(g, "## %s\n\n", file.Desc.Path())
	if pkg := string(file.Desc.Package()); pkg != "" {
		fmt.Fprintf(g, "- Package: `%s`\n\n", pkg)
	}

	if len(file.Services) > 0 {
		fmt.Fprintln(g, "### Services\n")
		for _, service := range file.Services {
			fmt.Fprintf(g, "#### %s\n\n", service.GoName)
			if desc := commentText(service.Comments.Leading); desc != "" {
				fmt.Fprintf(g, "%s\n\n", desc)
			}
			if len(service.Methods) == 0 {
				fmt.Fprintln(g, "_RPC は定義されていません。_\n")
				continue
			}

			fmt.Fprintln(g, "| RPC | リクエスト | レスポンス | ストリーミング |\n| --- | --- | --- | --- |")
			for _, method := range service.Methods {
				fmt.Fprintf(
					g,
					"| %s | `%s` | `%s` | %s |\n",
					method.GoName,
					messageIdent(method.Input),
					messageIdent(method.Output),
					streamingLabel(method),
				)
			}
			fmt.Fprintln(g)
		}
	}

	if len(file.Messages) > 0 {
		fmt.Fprintln(g, "### Messages\n")
		for _, message := range file.Messages {
			writeMessage(g, message, 4)
		}
	}

	if len(file.Enums) > 0 {
		fmt.Fprintln(g, "### Enums\n")
		for _, enum := range file.Enums {
			writeEnum(g, enum, 4)
		}
	}
}

func writeMessage(g *protogen.GeneratedFile, message *protogen.Message, headingLevel int) {
	if isMapEntry(message) {
		return
	}
	fmt.Fprintf(g, "%s %s\n\n", strings.Repeat("#", headingLevel), message.GoIdent.GoName)
	if desc := commentText(message.Comments.Leading); desc != "" {
		fmt.Fprintf(g, "%s\n\n", desc)
	}

	if len(message.Fields) > 0 {
		fmt.Fprintln(g, "| フィールド | 型 | ラベル | 説明 |\n| --- | --- | --- | --- |")
		for _, field := range message.Fields {
			fmt.Fprintf(
				g,
				"| %s | `%s` | %s | %s |\n",
				field.GoName,
				fieldType(field),
				fieldLabel(field),
				tableText(commentText(field.Comments.Leading)),
			)
		}
		fmt.Fprintln(g)
	} else {
		fmt.Fprintln(g, "_フィールドは定義されていません。_\n")
	}

	if len(message.Enums) > 0 {
		for _, enum := range message.Enums {
			writeEnum(g, enum, min(6, headingLevel+1))
		}
	}
	if len(message.Messages) > 0 {
		for _, nested := range message.Messages {
			writeMessage(g, nested, min(6, headingLevel+1))
		}
	}
}

func writeEnum(g *protogen.GeneratedFile, enum *protogen.Enum, headingLevel int) {
	fmt.Fprintf(g, "%s %s\n\n", strings.Repeat("#", headingLevel), enum.GoIdent.GoName)
	if desc := commentText(enum.Comments.Leading); desc != "" {
		fmt.Fprintf(g, "%s\n\n", desc)
	}
	if len(enum.Values) == 0 {
		fmt.Fprintln(g, "_値は定義されていません。_\n")
		return
	}

	fmt.Fprintln(g, "| 定数 | 番号 | 説明 |\n| --- | --- | --- |")
	for _, value := range enum.Values {
		fmt.Fprintf(
			g,
			"| %s | %d | %s |\n",
			value.GoIdent.GoName,
			value.Desc.Number(),
			tableText(commentText(value.Comments.Leading)),
		)
	}
	fmt.Fprintln(g)
}

func fieldType(field *protogen.Field) string {
	switch field.Desc.Kind() {
	case protoreflect.EnumKind:
		if field.Enum != nil {
			return string(field.Enum.Desc.FullName())
		}
	case protoreflect.MessageKind, protoreflect.GroupKind:
		if field.Message != nil {
			return string(field.Message.Desc.FullName())
		}
	}
	if field.Desc.IsMap() {
		return fmt.Sprintf("map<%s, %s>", field.Desc.MapKey().Kind(), field.Desc.MapValue().Kind())
	}
	return field.Desc.Kind().String()
}

func fieldLabel(field *protogen.Field) string {
	if field.Desc.IsMap() {
		return "map"
	}
	switch field.Desc.Cardinality() {
	case protoreflect.Repeated:
		return "repeated"
	case protoreflect.Required:
		return "required"
	default:
		if field.Oneof != nil && !field.Oneof.Desc.IsSynthetic() {
			return fmt.Sprintf("oneof %s", field.Oneof.GoIdent.GoName)
		}
		if field.Desc.HasOptionalKeyword() {
			return "optional"
		}
	}
	return "optional"
}

func streamingLabel(method *protogen.Method) string {
	clientStreaming := method.Desc.IsStreamingClient()
	serverStreaming := method.Desc.IsStreamingServer()
	switch {
	case clientStreaming && serverStreaming:
		return "双方向"
	case clientStreaming:
		return "クライアントストリーム"
	case serverStreaming:
		return "サーバーストリーム"
	default:
		return "-"
	}
}

func messageIdent(message *protogen.Message) string {
	if message == nil {
		return "-"
	}
	return string(message.Desc.FullName())
}

func commentText(comment protogen.Comments) string {
	trimmed := strings.TrimSpace(string(comment))
	if trimmed == "" {
		return ""
	}
	return collapseWhitespace(trimmed)
}

func tableText(value string) string {
	if value == "" {
		return "-"
	}
	value = strings.ReplaceAll(value, "|", "&#124;")
	return value
}

func collapseWhitespace(input string) string {
	if input == "" {
		return ""
	}
	fields := strings.Fields(input)
	return strings.Join(fields, " ")
}

func isMapEntry(message *protogen.Message) bool {
	if message.Desc.Options() == nil {
		return false
	}
	opts, ok := message.Desc.Options().(*descriptorpb.MessageOptions)
	if !ok {
		return false
	}
	return opts.GetMapEntry()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
