package pause

type Pause struct{
	SinglePauses []*SinglePause
}

type SinglePause struct{
	Pair string
	FromTime string
	ToTime string
	
}