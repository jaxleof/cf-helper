package mashup

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/jaxleof/uispinner"
)

func CloneContest(title string, id string, duration string) {
	cj := uispinner.New()
	cj.Start()
	defer cj.Stop()
	login := cj.AddSpinner(spinner.CharSets[34], 100*time.Millisecond).SetPrefix("cloning").SetComplete("clone complete")
	_, err := me.R().SetFormData(map[string]string{
		"action":                 "saveMashup",
		"isCloneContest":         "true",
		"parentContestIdAndName": id,
		"parentContestId":        id,
		"contestName":            title,
		"contestDuration":        duration,
		"problemsJson":           "[]",
		"csrf_token":             csrf,
	}).Post("https://codeforces.com/data/mashup")
	if err != nil {
		log.Fatal(err)
	}
	login.Done()
}

func QueryProbelmId(problem string) (string, error) {
	var info ProblemInfos
	_, err := me.R().SetFormData(map[string]string{
		"action":                      "problemQuery",
		"problemQuery":                problem,
		"previouslyAddedProblemCount": "0",
		"csrf_token":                  csrf,
	}).SetResult(&info).Post("https://codeforces.ml/data/mashup")
	if err != nil {
		log.Fatal(err)
	}
	if len(info.Problems) == 0 {
		log.Fatal(errors.New(problem + " isn't exist"))
	}
	return info.Problems[0].Id, nil
}

type ProblemInfo struct {
	EnglishName   string   `json:"englishName"`
	Id            string   `json:"id"`
	LocalizedName string   `json:"localizedName"`
	Rating        int      `json:"rating"`
	RussianName   string   `json:"russianName"`
	SolutionsUrl  string   `json:"solutionsUrl"`
	SolvedCount   int      `json:"solvedCount"`
	StatementUrl  string   `json:"statementUrl"`
	Tags          []string `json:"tags"`
}
type ProblemInfos struct {
	Problems []ProblemInfo `json:"problems"`
	Success  string        `json:"success"`
}

func CreateContest(title string, duration string, problems []string) {
	cj := uispinner.New()
	login := cj.AddSpinner(spinner.CharSets[34], 100*time.Millisecond).SetPrefix("contest creating").SetComplete("contest create complete")
	cj.Start()
	var problemsJson = make([]ProblemJson, len(problems))
	var group = new(sync.WaitGroup)
	var lock = new(sync.RWMutex)
	for i := 0; i < len(problems); i++ {
		group.Add(1)
		x := cj.AddSpinner(spinner.CharSets[34], 100*time.Millisecond).SetPrefix(problems[i] + " problem clawing").SetComplete(problems[i] + " claw complete")
		go func(index int) {
			id, err := QueryProbelmId(problems[index])
			if err != nil {
				cj.Stop()
				log.Fatalln(err)
				return
			}
			lock.Lock()
			problemsJson[index].Id = id
			problemsJson[index].Index = string(rune(('A' + index)))
			lock.Unlock()
			group.Done()
			x.Done()
		}(i)
	}
	group.Wait()
	data, err := json.Marshal(problemsJson)
	if err != nil {
		log.Fatal(err)
	}
	_, err = me.R().SetFormData(map[string]string{
		"action":                 "saveMashup",
		"isCloneContest":         "false",
		"parentContestIdAndName": "",
		"parentContestId":        "",
		"contestName":            title,
		"contestDuration":        duration,
		"problemsJson":           string(data),
		"csrf_token":             csrf,
	}).Post("https://codeforces.ml/data/mashup")
	if err != nil {
		log.Fatal(err)
	}
	login.Done()
	cj.Stop()
}

type ProblemJson struct {
	Id    string `json:"id"`
	Index string `json:"index"`
}
