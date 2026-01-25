package cli

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArgParser_Parse(t *testing.T) {
	const taskName = "some-task"
	remainingArgs := []string{"remaining", "--flags", "and", "[args]"}

	type FlagCase struct {
		flag          *Flag
		tokens        []string
		expectedValue any
	}

	flagCasesBuilder := func() map[string]FlagCase {
		return map[string]FlagCase{
			"string long": {
				flag: &Flag{
					Name:         "flag-a",
					Type:         FlagTypeString,
					DefaultValue: "",
				},
				tokens:        []string{"--flag-a", "flag-a-value"},
				expectedValue: "flag-a-value",
			},
			"string alias": {
				flag: &Flag{
					Name:         "flag-b",
					Type:         FlagTypeString,
					Aliases:      []string{"flag-b-alias"},
					DefaultValue: "",
				},
				tokens:        []string{"--flag-b-alias", "flag-b-value"},
				expectedValue: "flag-b-value",
			},
			"string short": {
				flag: &Flag{
					Name:         "flag-c",
					Type:         FlagTypeString,
					Shorts:       []rune{'c'},
					DefaultValue: "",
				},
				tokens:        []string{"-c", "flag-c-value"},
				expectedValue: "flag-c-value",
			},
			"string long equals": {
				flag: &Flag{
					Name:         "flag-d",
					Type:         FlagTypeString,
					DefaultValue: "",
				},
				tokens:        []string{"--flag-d=flag-d-value"},
				expectedValue: "flag-d-value",
			},
			"string alias equals": {
				flag: &Flag{
					Name:         "flag-e",
					Type:         FlagTypeString,
					Aliases:      []string{"flag-e-alias"},
					DefaultValue: "",
				},
				tokens:        []string{"--flag-e-alias=flag-e-value"},
				expectedValue: "flag-e-value",
			},
			"string short equals": {
				flag: &Flag{
					Name:         "flag-f",
					Type:         FlagTypeString,
					Shorts:       []rune{'f'},
					DefaultValue: "",
				},
				tokens:        []string{"-f=flag-f-value"},
				expectedValue: "flag-f-value",
			},
			"bool long": {
				flag: &Flag{
					Name:         "flag-g",
					Type:         FlagTypeBool,
					DefaultValue: false,
				},
				tokens:        []string{"--flag-g"},
				expectedValue: true,
			},
			"bool alias": {
				flag: &Flag{
					Name:         "flag-h",
					Type:         FlagTypeBool,
					Aliases:      []string{"flag-h-alias"},
					DefaultValue: false,
				},
				tokens:        []string{"--flag-h-alias"},
				expectedValue: true,
			},
			"bool short": {
				flag: &Flag{
					Name:         "flag-i",
					Type:         FlagTypeBool,
					Shorts:       []rune{'i'},
					DefaultValue: false,
				},
				tokens:        []string{"-i"},
				expectedValue: true,
			},
			"bool long equals false": {
				flag: &Flag{
					Name:         "flag-j",
					Type:         FlagTypeBool,
					DefaultValue: false,
				},
				tokens:        []string{"--flag-j=false"},
				expectedValue: false,
			},
			"bool long equals true": {
				flag: &Flag{
					Name:         "flag-k",
					Type:         FlagTypeBool,
					DefaultValue: false,
				},
				tokens:        []string{"--flag-k=true"},
				expectedValue: true,
			},
			"bool alias equals false": {
				flag: &Flag{
					Name:         "flag-l",
					Type:         FlagTypeBool,
					Aliases:      []string{"flag-l-alias"},
					DefaultValue: false,
				},
				tokens:        []string{"--flag-l-alias=0"},
				expectedValue: false,
			},
			"bool alias equals true": {
				flag: &Flag{
					Name:         "flag-m",
					Type:         FlagTypeBool,
					Aliases:      []string{"flag-m-alias"},
					DefaultValue: false,
				},
				tokens:        []string{"--flag-m-alias=1"},
				expectedValue: true,
			},
			"bool short equals false": {
				flag: &Flag{
					Name:         "flag-n",
					Type:         FlagTypeBool,
					Shorts:       []rune{'n'},
					DefaultValue: false,
				},
				tokens:        []string{"-n=f"},
				expectedValue: false,
			},
			"bool short equals true": {
				flag: &Flag{
					Name:         "flag-o",
					Type:         FlagTypeBool,
					Shorts:       []rune{'o'},
					DefaultValue: false,
				},
				tokens:        []string{"-o=TRUE"},
				expectedValue: true,
			},
			"bool long invert": {
				flag: &Flag{
					Name:         "flag-p",
					Type:         FlagTypeBool,
					DefaultValue: true,
				},
				tokens:        []string{"--flag-p"},
				expectedValue: false,
			},
		}
	}

	t.Run("kitchen sink", func(t *testing.T) {
		flagCases := flagCasesBuilder()

		flags := lo.Map(lo.Values(flagCases), func(fc FlagCase, _ int) *Flag { return fc.flag })
		tokens := lo.FlatMap(lo.Values(flagCases), func(fc FlagCase, _ int) []string { return fc.tokens })

		parser := NewParser("v0.1.0", append(append(tokens, taskName), remainingArgs...), flags...)
		require.NoError(t, parser.Parse())

		assert.Equal(t, taskName, parser.taskName)
		assert.Equal(t, remainingArgs, parser.RemainingArgs())

		for _, fc := range flagCases {
			assert.Equal(t, fc.expectedValue, fc.flag.Value())
		}
	})

	t.Run("individual", func(t *testing.T) {
		flagCases := flagCasesBuilder()

		for name, fc := range flagCases {
			t.Run(name, func(t *testing.T) {
				parser := NewParser("v0.1.0", append(append(fc.tokens, taskName), remainingArgs...), fc.flag)
				require.NoError(t, parser.Parse())

				assert.Equal(t, fc.expectedValue, fc.flag.Value())
			})
		}
	})
}
