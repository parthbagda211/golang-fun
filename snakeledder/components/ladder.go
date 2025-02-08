package components

type Ladder struct{
	Start int  
	End int 
}

func NewLeadder(start,end int ) *Ladder {
	return &Ladder{Start: start,End: end}
}