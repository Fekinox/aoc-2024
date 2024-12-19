package main

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/Fekinox/aoc-2024/util"
)

func Problem13NoRegex(in io.Reader, out io.Writer) {
	var silver, gold int64

	groups := util.ReadNewlineSeparatedGroups(in)

	for _, g := range groups {
		aTokens := strings.Split(g[0], " ")
		bTokens := strings.Split(g[1], " ")
		prizeTokens := strings.Split(g[2], " ")

		axt, ayt := aTokens[2], aTokens[3]
		bxt, byt := bTokens[2], bTokens[3]
		pxt, pyt := prizeTokens[1], prizeTokens[2]

		ax, ay := util.MustParseInt(axt[2:len(axt)-1]), util.MustParseInt(ayt[2:])
		bx, by := util.MustParseInt(bxt[2:len(bxt)-1]), util.MustParseInt(byt[2:])
		px, py := util.MustParseInt(pxt[2:len(pxt)-1]), util.MustParseInt(pyt[2:])
		gpx, gpy := px+10000000000000, py+10000000000000

		det := ax*by - bx*ay

		sa, sb := (by*px-bx*py)/det, (-ay*px+ax*py)/det
		ga, gb := (by*gpx-bx*gpy)/det, (-ay*gpx+ax*gpy)/det

		if sa*ax+sb*bx == px && sa*ay+sb*by == py {
			silver += sa*3 + sb
		}

		if ga*ax+gb*bx == gpx && ga*ay+gb*by == gpy {
			gold += ga*3 + gb
		}
	}

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem13Regex(in io.Reader, out io.Writer) {
	var silver, gold int64

	groups := util.ReadNewlineSeparatedGroups(in)

	buttonRegex := regexp.MustCompile("(X|Y)\\+(\\d+)")
	prizeRegex := regexp.MustCompile("(X|Y)=(\\d+)")

	for _, g := range groups {
		aTokens := buttonRegex.FindAllStringSubmatch(g[0], -1)
		bTokens := buttonRegex.FindAllStringSubmatch(g[1], -1)
		prizeTokens := prizeRegex.FindAllStringSubmatch(g[2], -1)

		ax, ay := util.MustParseInt(aTokens[0][2]), util.MustParseInt(aTokens[1][2])
		bx, by := util.MustParseInt(bTokens[0][2]), util.MustParseInt(bTokens[1][2])
		px, py := util.MustParseInt(prizeTokens[0][2]), util.MustParseInt(prizeTokens[1][2])
		gpx, gpy := px+10000000000000, py+10000000000000

		det := ax*by - bx*ay

		sa, sb := (by*px-bx*py)/det, (-ay*px+ax*py)/det
		ga, gb := (by*gpx-bx*gpy)/det, (-ay*gpx+ax*gpy)/det

		if sa*ax+sb*bx == px && sa*ay+sb*by == py {
			silver += sa*3 + sb
		}

		if ga*ax+gb*bx == gpx && ga*ay+gb*by == gpy {
			gold += ga*3 + gb
		}
	}

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem13(in io.Reader, out io.Writer) {
	Problem13NoRegex(in, out)
}
