/*
Test of class setter/getter for rules engine
*/

package main

import (
	"fmt"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	"log"
)

const (
	SourceTypeLdap   = "ldap"
	SourceTypeVMWare = "vmware"
)

func NewCmdbCiServer() *CmdbCiServer {
	obj := &CmdbCiServer{}
	obj.Hostname.configurationItem = make(map[string]string)
	return obj
}

type CmdbCiServer struct {
	Hostname Mutator[string] `json:"hostname,omitempty"`
	Owner    Mutator[string] `json:"owner,omitempty"`
	Version  Mutator[int]
	Active   Mutator[bool]
}

type Mutator[T any] struct {
	value             T
	configurationItem map[string]T
	ldap              T
	vmWare            T
	err               T
}

func (mt *Mutator[T]) SetLdap(val T) *Mutator[T] {
	mt.configurationItem[SourceTypeLdap] = val
	return mt
}

func (mt *Mutator[T]) UseLdap() {
	mt.value = mt.configurationItem[SourceTypeLdap]
}

func (mt *Mutator[T]) SetVMWare(val T) *Mutator[T] {
	mt.configurationItem[SourceTypeVMWare] = val
	return mt
}

func (mt *Mutator[T]) UseVMWare() {
	mt.value = mt.configurationItem[SourceTypeVMWare]
}

func (mt *Mutator[T]) GetValue() *Mutator[T] {
	return &Mutator[T]{value: mt.value}
}

// Helper methods

func (mt *Mutator[T]) String() string {
	return fmt.Sprintf("%v", mt.value)
}

func main() {
	test := NewCmdbCiServer()

	test.Hostname.SetLdap("appsrv01.domain.com")
	test.Hostname.SetVMWare("server1 (used as an app server)")

	fmt.Printf("value is now: %s\n", test.Hostname.GetValue())
	test.Hostname.UseVMWare()
	fmt.Printf("value is now: %s\n", test.Hostname.GetValue())
	test.Hostname.UseLdap()
	fmt.Printf("value is now: %s\n", test.Hostname.GetValue())

	dataCtx := ast.NewDataContext()
	err := dataCtx.Add("MF", test)
	if err != nil {
		log.Printf("%v", err)
	}
	knowledgeLibrary := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(knowledgeLibrary)

	// lets prepare a rule definition
	drls := `
	rule CheckValues "Check the default values" salience 10 {
    when 
		true
    then
        MF.Hostname.UseVMWare();
        Retract("CheckValues");
	}`

	// Add the rule definition above into the library and name it 'TutorialRules'  version '0.0.1'
	bs := pkg.NewBytesResource([]byte(drls))
	err = ruleBuilder.BuildRuleFromResource("TutorialRules", "0.0.1", bs)
	if err != nil {
		log.Printf("%v", err)
	}

	knowledgeBase := knowledgeLibrary.NewKnowledgeBaseInstance("TutorialRules", "0.0.1")

	engine := engine.NewGruleEngine()
	err = engine.Execute(dataCtx, knowledgeBase)
	if err != nil {
		log.Panic(err)
	}

	log.Print(test.Hostname.GetValue())
}
