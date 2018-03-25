package storage_test

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	f "github.com/simonlissack/footballfixtures/ffconfig"
	"github.com/simonlissack/footballfixtures/model"
	. "github.com/simonlissack/footballfixtures/storage"
)

var _ = Describe("localTeamsCacheClient", func() {
	const storedTeamsCache = "teamsCache_test.json"
	var (
		config     f.FFConfiguration
		teamsCache TeamsCacheClient
		teams      []model.Team
		err        error
	)

	Describe("When loading a file", func() {

		JustBeforeEach(func() {
			teamsCache = NewLocalTeamsCache(config)
			teams, err = teamsCache.LoadTeams()
		})

		Context("which exists", func() {
			BeforeEach(func() {
				config = f.FFConfiguration{TeamsFile: storedTeamsCache}
			})

			It("should load the teams", func() {
				Expect(len(teams)).To(Equal(2))
			})

			It("should not throw an error", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("which does not exists", func() {
			BeforeEach(func() {
				config = f.FFConfiguration{TeamsFile: "teamsCache_test_not_found.json"}
			})

			It("should load the teams", func() {
				Expect(teams).To(BeNil())
			})

			It("should not throw an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(ContainSubstring("GetFileAttributesEx teamsCache_test_not_found.json: The system cannot find the file specified."))
			})
		})
	})

	Describe("When storing teams", func() {
		const teamCacheSaveLoc = "teamsCache_test_store.json"
		var (
			teams       []model.Team
			loadedTeams []model.Team
		)

		BeforeEach(func() {
			config = f.FFConfiguration{TeamsFile: "teamsCache_test.json"}
			teamsCache = NewLocalTeamsCache(config)
			teams, _ = teamsCache.LoadTeams()
		})

		JustBeforeEach(func() {
			config = f.FFConfiguration{TeamsFile: teamCacheSaveLoc}
			teamsCache = NewLocalTeamsCache(config)
			err = teamsCache.SaveTeams(teams)
			loadedTeams, _ = teamsCache.LoadTeams()
		})

		It("should store the teams file", func() {
			_, err := os.Stat(teamCacheSaveLoc)
			Expect(os.IsNotExist(err)).To(Equal(false), fmt.Sprintf("File '%s' does not exist", config.TeamsFile))
		})

		It("should load the teams", func() {
			Expect(len(teams)).To(Equal(2))
		})

		It("should not throw an error", func() {
			Expect(err).To(BeNil())
		})

		AfterEach(func() {
			os.Remove(config.TeamsFile)
		})
	})
})
