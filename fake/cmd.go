package fake

import (
	"log"
	"net/url"

	"github.com/SecurityBrewery/catalyst/generated/model"

	"github.com/SecurityBrewery/catalystctl/client"
)

type Cmd struct {
	Token string   `help:"Token to use for authentication."`
	URL   *url.URL `help:"URL of the Catalyst server."`
}

func (c *Cmd) Run() error {
	catalystClient, err := client.New(c.URL, c.Token)
	if err != nil {
		return err
	}

	g := &Generator{client: catalystClient}

	if err := g.userDummyData(); err != nil {
		return err
	}
	if err := g.ticketDummyData(); err != nil {
		return err
	}
	if err := g.dashboardDummyData(); err != nil {
		return err
	}

	return nil
}

type Generator struct {
	client *client.CatalystClient
}

var users = []*model.UserForm{
	{ID: "alice", Blocked: false, Roles: []string{"analyst"}},
	{ID: "bob", Blocked: false, Roles: []string{"analyst"}},
	{ID: "carol", Blocked: false, Roles: []string{"analyst"}},
	{ID: "dave", Blocked: false, Roles: []string{"analyst"}},
	{ID: "eve", Blocked: false, Roles: []string{"admin"}},
}

// var settings = []*models.Setting{
// 	{Email: swag.String("alice@example.com"), Name: swag.String("Alice Alert Analyst"),, : "alice"},
// 	{Email: swag.String("bob@example.com"), Name: swag.String("Bob Incident Handler"), Username: "bob"},
// 	{Email: swag.String("carol@example.com"), Name: swag.String("Carol Forensicator"), Username: "carol"},
// 	{Email: swag.String("dave@example.com"), Name: swag.String("Dave Admin"), Username: "dave"},
// 	{Email: swag.String("eve@example.com"), Name: swag.String("Eve Team Lead"), Username: "eve"},
// }

func (g *Generator) dashboardDummyData() error {
	simple := &model.Dashboard{
		Name: "Simple",
		Widgets: []*model.Widget{
			{
				Aggregation: "type",
				Filter:      pointer("status == 'open'"),
				Name:        "Types",
				Type:        "pie",
				Width:       8,
			},
			{
				Aggregation: "owner",
				Filter:      pointer("status == 'open'"),
				Name:        "Owners",
				Type:        "bar",
				Width:       4,
			},
		},
	}

	log.Println("create dashboard")

	_, err := g.client.PostJSON("/dashboards", nil, simple)
	return err
}

func (g *Generator) userDummyData() error {
	for _, user := range users {
		log.Println("create user ", user.ID)
		_, _ = g.client.PostJSON("/users", nil, user)
	}
	return nil
}

func (g *Generator) ticketDummyData() error {
	if err := g.createTickets(10_000, fakeIncident); err != nil {
		return err
	}
	if err := g.createTickets(200_000, fakeAlert); err != nil {
		return err
	}
	if err := g.createTickets(100, fakeCustomTicketInvestigation); err != nil {
		return err
	}
	if err := g.createTickets(240, fakeCustomTicketHunt); err != nil {
		return err
	}

	return nil
}

func (g *Generator) createTickets(count int, createFunc func() *model.TicketForm) error {
	log.Println("create ticket")
	var tickets []*model.TicketForm
	for j := 0; j < count; j++ {
		tickets = append(tickets, createFunc())

		if len(tickets) > 100 {
			if _, err := g.client.PostJSON("/tickets/batch", nil, tickets); err != nil {
				return err
			}
			tickets = nil
		}
	}
	if len(tickets) > 0 {
		if _, err := g.client.PostJSON("/tickets/batch", nil, tickets); err != nil {
			return err
		}
	}
	return nil
}
