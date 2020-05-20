package cmd

import (
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/liclac/sharlayan/calibre"
	"github.com/liclac/sharlayan/render"
	"github.com/liclac/sharlayan/tree"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build a static website",
	Long: `Build a static website.

This command requires you to specify (at least one) style of path to use.

## By Name - 'books/Equal Rites', 'authors/Terry Pratchett'.
$ sharlayan build -n
$ sharlayan build --build.by-name.enable

This style is similar to Calibre's ('Terry Pratchett/The Last Elephant (41)'),
with some more flexibility, and is ideal for browsing your library by hand, or
for exposing it via a file sharing protocol like SFTP or Windows file sharing.

-------------------------------------------------------------------------------
Books
  Equal Rites
    metadata.json // { "id": 35, "title": "Equal Rites", ... }
    Files
      Terry Pratchett - Equal Rites.epub
  The Last Elephant
    metadata.json // { "id": 41, "title": "The Last Elephant", ... }
    Files
      Terry Pratchett - The Last Elephant.epub
Authors
  Terry Pratchett
    metadata.json // { "id": 4, "name": "Terry Pratchett", ... }
    Books
      Equal Rites -> ../../../Books/Equal Rites
      The Last Elephant -> ../../../Books/The Last Elephant
-------------------------------------------------------------------------------

## By ID - 'books/41', 'authors/4'.
$ sharlayan build -i
$ sharlayan build --build.by-id.enable

This style is ideal for using with a web server or similar, and produces paths
that don't change, even if things are moved or renamed. As a bonus, the paths
are much shorter than their named equivalent.

-------------------------------------------------------------------------------
books
  35
    metadata.json // { "id": 35, "title": "Equal Rites", ... }
    files
      Terry Pratchett - Equal Rites.epub
  41
    metadata.json // { "id": 41, "title": "The Last Elephant", ... }
    files
      Terry Pratchett - The Last Elephant.epub
authors
  4
    metadata.json // { "id": 4, "name": "Terry Pratchett", ... }
    books
      Equal Rites -> ../../../books/35
      The Last Elephant -> ../../../books/41
-------------------------------------------------------------------------------

## Why not both?
$ sharlayan build -i -n
$ sharlayan build --build.by-id.enable --build.by-name.enable

You can just enable both styles, to create a "stacked layout" of both:
- /books/35
- /books/Equal Rites -> ./35
- /authors/4
- /authors/Terry Pratchett -> ./4

The catch is that stacks get quite cluttered. This is when prefixes come in:

$ sharlayan build -I by-id -N by-name
$ sharlayan build --build.by-id.prefix=by-id --build.by-name.prefix=by-name

Passing '-I'/'--build.by-id.prefix' implies '-i'/'--build.by-id.enable'.
Passing '-N'/'--build.by-name.prefix' implies '-n'/'--build.by-name.enable'.

- /by-id/books/35
- /by-id/authors/4
- /by-name/Books/Equal Rites -> ../../by-id/books/35
- /by-name/Authors/Terry Pratchett -> ../../by-id/authors/4

You can also put only one under a prefix, eg. '-n -I _id' to create a named
tree with an ID tree under '/_id'.

## Deduplication and LinkIDs

Sharlayan's output is always deduplicated - when using both naming schemes,
'/by-name/Books/Wyrd Sisters' will turn into a symlink to '/by-id/books/35'.
If this isn't desirable, you can just build both trees separately.

The reason for this is that all references are resolved using "Link IDs", which
means that an HTML template or symlink can just link to "books:35", which is
resolved to an actual path at runtime.

Link IDs are unique, and if a second instance is added, it turns into a symlink
to the first one. The named tree is always the latter one added, simply because
symlinks tend to be much less problematic to a human user in a file explorer,
than to eg. a web server, some of which refuse to follow symlinks by default.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		meta, err := calibre.Read(cfg.Library)
		if err != nil {
			return err
		}
		fs := osfs.New(cfg.Build.Out)
		return tree.Render(fs, false, "/", render.Root(cfg, meta))
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringP("build.out", "o", "out", "path to output")

	buildCmd.Flags().BoolP("build.by-name.enable", "n", false, "use names in paths (implied by -N)")
	buildCmd.Flags().StringP("build.by-name.prefix", "N", "", "prefix for named paths")
	buildCmd.Flags().BoolP("build.by-id.enable", "i", false, "use IDs in paths (implied by -I)")
	buildCmd.Flags().StringP("build.by-id.prefix", "I", "", "prefix for ID paths")

	viper.BindPFlags(buildCmd.Flags())
}
