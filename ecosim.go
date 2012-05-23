
package ecosim

import (
	"time"
	. "../gomarket"
)

type Engine struct {
	actors map[Actor]bool
	market *Market
}
func NewEngine() *Engine {
	return &Engine{make(map[Actor]bool), NewMarket()}
}
func (e *Engine) Add(a Actor) {
	e.actors[a] = true
	e.market.Add(a)
}
func (e *Engine) Del(a Actor) {
	delete(e.actors, a)
	e.market.Del(a)
}
func (e *Engine) Run(t time.Duration) {
	for actor,_ := range e.actors {
		profitCalculator := NewProfitCalculator(e, actor, t)
		best_profit := Profit(nil)
		next_best_profit := Profit(nil)
		for _,process := range actor.Processes() {
			profit := profitCalculator.processProfit(process, t)
			if profit > best_profit {
				next_best_profit = best_profit
				best_process, best_profit, best_output = process, profit, output
			} else if profit > next_best_profit {
				next_best_profit = profit
			}
		}
		if best_output != nil {
			actor.Update(&Update{best_output, t, next_best_profit})
		}
	}
	e.market.Trade()
}


/*
 * The immediate and eventual output of a process.
 * Eventual output is defined as actual output + the fraction of not-yet-produced output for the next cycle.
 */
type Output struct {
	Process *Process
	Immediate Resources
	Eventual Resources
}
func (o *Output) MergeIn(other *Output) {
	o.Immediate.MergeIn(other.Immediate)
	o.Eventual.MergeIn(other.Eventual)
}
func (o *Output) Profit(market *Market) {
	return &Profit{o, market.Value(o.Immediate), market.Value(o.Eventual)}
}

type Update struct {
	Profit *profit
	Time time.Duration
}

/*
 * The Output of a process along with its immedate and eventual profit (or, the avoidance costs along with their loss)
 */
type Profit struct {
	output *Output
	Immediate float64
	Eventual float64
}
func (p *Profit) MergeIn(o *Profit) {
	p.output.MergeIn(o.output)
	p.Immediate += o.Immediate
	p.Eventual += p.Eventual
}

type Process interface {
	/*
	 * The results when running this process for a time.
	 */
	Run(time.Duration) *Output
	/*
	 * The results when avoiding this process for a time.
	 */
	Avoid(time.Duration) *Output
}

type Actor interface {
	Trader
	Carrier
	Processes() []Process
	Update(*Update)
}

type Skill interface{}

type Skills map[Skill]float64

type StandardActor struct {
	StandardTrader
	processes []Process
	skills Skills
	resources Resources
}
func (s *StandardActor) AddProcess(factory *StandardProcessFactory) {
	s.processes = append(s.processes, factory.produce(s))
}
func (s *StandardActor) Processes() []Process {
	return s.processes
}
func (s *StandardActor) Update(update *Update) {
	fmt.Println(s,"updated with",update)
}
