package tests

type A struct{}

func (a A) Wat() {

}

func (a A) Wot() {
	a.Wat()
}
