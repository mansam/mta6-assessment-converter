/*
Convert Pathfinder assessments to Konveyor Hub assessments,
assuming that the legacy questionnaire has been seeded at ID 1.

Accepts a JSON document containing a list of Pathfinder assessments
and outputs a list of Hub assessments to standard out.
*/

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/konveyor/tackle2-hub/api"
	"github.com/konveyor/tackle2-hub/assessment"
)

type PathfinderAssessment struct {
	ID                uint
	ApplicationID     uint
	Status            string
	Stakeholders      []uint
	StakeholderGroups []uint
	Questionnaire     Questionnaire
}

type Questionnaire struct {
	Categories []Category
	Title      string
	Language   string
}

type Category struct {
	ID        uint
	Order     uint
	Title     string
	Questions []Question
	Comment   string
}

type Question struct {
	ID          uint
	Order       uint
	Question    string
	Description string
	Options     []Option
}

type Option struct {
	ID      uint
	Order   uint
	Option  string
	Checked bool
	Risk    string
}

type Mappable[T any, S any] func(T) S

func Map[T any, S any](f Mappable[T, S], in []T) (out []S) {
	for _, each := range in {
		out = append(out, f(each))
	}
	return
}

func Ref(id uint) (ref api.Ref) {
	ref.ID = id
	return
}

func convert(p PathfinderAssessment) (a api.Assessment) {
	a.Questionnaire = api.Ref{ID: 1}
	a.Application = &api.Ref{ID: p.ApplicationID}
	a.Sections = Map(convertCategory, p.Questionnaire.Categories)
	a.Stakeholders = Map(Ref, p.Stakeholders)
	a.StakeholderGroups = Map(Ref, p.StakeholderGroups)
	return
}

func convertCategory(c Category) (section assessment.Section) {
	section.Name = c.Title
	section.Order = &c.Order
	section.Comment = c.Comment
	section.Questions = Map(convertQuestion, c.Questions)
	return
}

func convertQuestion(q Question) (question assessment.Question) {
	question.Text = q.Question
	question.Explanation = q.Description
	question.Order = &q.Order
	question.Answers = Map(convertOption, q.Options)
	return
}

func convertOption(o Option) (answer assessment.Answer) {
	answer.Text = o.Option
	answer.Selected = o.Checked
	answer.Order = &o.Order
	risk := strings.ToLower(o.Risk)
	switch risk {
	case "amber":
		answer.Risk = "yellow"
	default:
		answer.Risk = risk
	}
	return
}

func usage() {
	fmt.Printf("usage: %s path-to-assessments.json\n", os.Args[0])
	os.Exit(1)
}

func main() {
	if len(os.Args) != 2 {
		usage()
	}
	bytes, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("error reading file '%s': %s\n", os.Args[1], err)
		os.Exit(1)
	}

	unconverted := []PathfinderAssessment{}
	err = json.Unmarshal(bytes, &unconverted)
	if err != nil {
		fmt.Printf("error decoding assessment json: %s", err)
		os.Exit(1)
	}
	var converted []api.Assessment
	converted = Map(convert, unconverted)
	output, err := json.MarshalIndent(converted, "", "    ")
	if err != nil {
		fmt.Printf("error marshalling output json: %s", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", output)
}
