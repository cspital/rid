package main

import (
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rs/xid"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
)

var out = os.Stdout

var root = &cobra.Command{
	Use:   "rid",
	Short: "Generate some ids",
}

var num int

func init() {
	root.PersistentFlags().IntVarP(&num, "num", "n", 1, "number of ids to make")
	u := &cobra.Command{
		Use:   "u",
		Short: "UUID",
		Run: func(cmd *cobra.Command, args []string) {
			g := uu{}
			generate(g, num)
		},
	}
	x := &cobra.Command{
		Use:   "x",
		Short: "Xid",
		Run: func(cmd *cobra.Command, args []string) {
			g := xi{}
			generate(g, num)
		},
	}
	l := &cobra.Command{
		Use:   "l",
		Short: "Ulid",
		Run: func(cmd *cobra.Command, args []string) {
			g := ul{
				entropy: ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0),
			}
			generate(g, num)
		},
	}
	root.AddCommand(u)
	root.AddCommand(x)
	root.AddCommand(l)
}

func main() {
	root.Execute()
}

func generate(g generator, n int) {
	for i := 0; i < n; i++ {
		out.WriteString(g.Generate())
		if i < n {
			out.WriteString("\n")
		}
	}
}

type generator interface {
	Generate() string
}

type uu struct{}

func (u uu) Generate() string {
	return uuid.NewV4().String()
}

type xi struct{}

func (x xi) Generate() string {
	return xid.New().String()
}

type ul struct {
	entropy io.Reader
}

func (u ul) Generate() string {
	t := time.Now()
	return ulid.MustNew(ulid.Timestamp(t), u.entropy).String()
}
