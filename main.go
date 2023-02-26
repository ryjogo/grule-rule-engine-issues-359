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
	obj := &CmdbCiServer{
		Hostname: newMutator(),
		Owner:    newMutator(),
		Version:  newMutator(),
		Active:   newMutator(),
	}
	return obj
}

type CmdbCiServer struct {
	Hostname *mutator
	Owner    *mutator
	Version  *mutator
	Active   *mutator
}

type mutator struct {
	Value             interface{}
	configurationItem map[string]interface{}
	ldap              interface{}
	vmWare            interface{}
	err               interface{}
}

func newMutator() *mutator {
	return &mutator{
		configurationItem: make(map[string]interface{}),
	}
}

func (mt *mutator) SetLdap(val interface{}) *mutator {
	mt.configurationItem[SourceTypeLdap] = val
	return mt
}

func (mt *mutator) UseLdap() {
	mt.Value = mt.configurationItem[SourceTypeLdap]
}

func (mt *mutator) SetVMWare(val interface{}) *mutator {
	mt.configurationItem[SourceTypeVMWare] = val
	return mt
}

func (mt *mutator) UseVMWare() {
	mt.Value = mt.configurationItem[SourceTypeVMWare]
}

func (mt *mutator) GetValue() interface{} {
	return mt.Value
}

// Helper methods

func (mt *mutator) String() string {
	return fmt.Sprintf("%v", mt.Value)
}

type CmdbCiServerWrapper struct {
	MF *CmdbCiServer
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
	err := dataCtx.Add("MF", &CmdbCiServerWrapper{test})
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
		Log("Value: " + MF.MF.Hostname.Value);
		MF.MF.Hostname.UseVMWare();
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
