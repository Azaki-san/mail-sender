package main

import (
	"example.com/mod/objects"
	"example.com/mod/objects/pipeline"
)

func initPipeline() *pipeline.Pipeline {
	pipelineInstance := &pipeline.Pipeline{}
	pipelineInstance.AddJob(&objects.WordsFilter{})
	pipelineInstance.AddJob(&objects.ScreamingFilter{})
	pipelineInstance.AddJob(&objects.PublishFilter{})
	return pipelineInstance
}
