package pipeline

import "example.com/mod/objects"

type Job struct {
	filter objects.Filter
	pipe   chan objects.Message
}

type Pipeline struct {
	jobs []Job
}

func (p *Pipeline) AddJob(filter objects.Filter) {
	p.jobs = append(p.jobs, Job{filter: filter, pipe: make(chan objects.Message)})
}

func (p *Pipeline) Process(msg objects.Message) objects.Message {
	pipelineResult := make(chan objects.Message)
	for i, job := range p.jobs {
		var channel chan objects.Message
		if i == len(p.jobs)-1 {
			channel = pipelineResult
		} else {
			channel = p.jobs[i+1].pipe
		}
		go func(channel chan objects.Message) {
			applied := job.filter.Apply(<-p.jobs[i].pipe)
			channel <- *applied
		}(channel)
	}
	p.jobs[0].pipe <- msg
	result := <-pipelineResult
	return result
}
