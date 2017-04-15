/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   main.go                                            :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: crenfrow <crenfrow@student.42.us>          +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2017/04/14 13:07:08 by crenfrow          #+#    #+#             */
/*   Updated: 2017/04/14 19:43:34 by crenfrow         ###   ########.fr       */
/*                                                                            */
/* ************************************************************************** */

package main

import (
	"crypto/tls"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	gomail "gopkg.in/gomail.v2"
)

// My custom email struct for hauling around the variables
type ftEmail struct {
	recipient *string
	recipPath *string
	subject   *string
	htmlPath  *string
	attPath   *string
}

// Reading and parsing of the CSV file into a string array given the path to the
// file
func parseCSV(path string) []string {

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	s := string(data)
	r := csv.NewReader(strings.NewReader(s))
	emailArr, err := r.ReadAll()
	if err != nil {
		panic(err)
	}

	return emailArr[0]
}

// Takes in our email struct and uses it to fill out the gomail object, then
// sends it to the recipent(s)
func sendMail(e ftEmail) {

	m := gomail.NewMessage()

	m.SetHeader("From", m.FormatAddress("crenfrow@student.42.us.org", "Chris Renfrow"))
	if *e.recipPath != "" {
		emailArr := parseCSV(*e.recipPath)
		m.SetHeader("To", emailArr...)
	} else {
		m.SetHeader("To", *e.recipient)
	}

	m.SetHeader("Subject", *e.subject)
	data, err := ioutil.ReadFile(*e.htmlPath)
	if err != nil {
		panic(err)
	}

	m.SetBody("text/html", string(data))
	if *e.attPath != "" {
		m.Attach(*e.attPath)
	} else {
		fmt.Println("No attachment!")
	}

	d := gomail.NewDialer("smtp.42.us.org", 25, "", "")
	// Essential line when using insecure
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

// Initializes flags for our CLI to use, parses those flags and assigns the
// values to our custom ftEmail struct of pointers
func parseFlags() *ftEmail {

	var to = flag.String("t", "",
		"Email address of recipient")
	var toList = flag.String("T", "",
		"A path to a CSV of recipients' emails")
	var subject = flag.String("s", "You've got mail!",
		"A friendly subject for your email")
	var message = flag.String("m", "default.html",
		"A valid html for the body of your message")
	var attachment = flag.String("a", "",
		"Direct path to attachment")

	flag.Parse()

	e := new(ftEmail)
	e.recipient = to
	e.recipPath = toList
	e.subject = subject
	e.htmlPath = message
	e.attPath = attachment

	return e
}

func main() {
	e := parseFlags()
	sendMail(*e)
}
