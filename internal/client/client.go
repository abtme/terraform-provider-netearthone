package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Client is the NetearthOne API client.
type Client struct {
	BaseURL    string
	AuthUserID int
	APIKey     string
	HTTPClient *http.Client
}

// NewClient creates a new API client.
func NewClient(baseURL string, authUserID int, apiKey string) *Client {
	return &Client{
		BaseURL:    strings.TrimRight(baseURL, "/"),
		AuthUserID: authUserID,
		APIKey:     apiKey,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// ---------------------------------------------------------------------------
// Shared response types
// ---------------------------------------------------------------------------

// actionResponse is the common response shape for modify/action operations.
type actionResponse struct {
	Status           string `json:"status"`
	Error            string `json:"error"`
	Description      string `json:"description"`
	EntityID         interface{} `json:"entityid"`
	ActionType       string `json:"actiontype"`
	ActionTypeDesc   string `json:"actiontypedesc"`
	EAQID            string `json:"eaqid"`
	ActionStatus     string `json:"actionstatus"`
	ActionStatusDesc string `json:"actionstatusdesc"`
}

// ---------------------------------------------------------------------------
// Domain details types
// ---------------------------------------------------------------------------

// DomainDetails holds domain information returned by the details endpoints.
type DomainDetails struct {
	OrderID    interface{} `json:"orderid"`
	DomainName string      `json:"domainname"`
	// Nameservers — returned as ns1..ns13
	NS1  string `json:"ns1"`
	NS2  string `json:"ns2"`
	NS3  string `json:"ns3"`
	NS4  string `json:"ns4"`
	NS5  string `json:"ns5"`
	NS6  string `json:"ns6"`
	NS7  string `json:"ns7"`
	NS8  string `json:"ns8"`
	NS9  string `json:"ns9"`
	NS10 string `json:"ns10"`
	NS11 string `json:"ns11"`
	NS12 string `json:"ns12"`
	NS13 string `json:"ns13"`
	// Privacy
	PrivacyProtected        interface{} `json:"privacyprotected"`
	PrivacyProtectedTilDate string      `json:"privacyprotectedtildate"`
	// Contacts
	RegistrantContactID interface{} `json:"registrantcontactid"`
	AdminContactID      interface{} `json:"admincontactid"`
	TechContactID       interface{} `json:"techcontactid"`
	BillingContactID    interface{} `json:"billingcontactid"`
	// Child nameservers — map of hostname -> []IPs
	CNS map[string]interface{} `json:"cns"`
}

// OrderIDInt returns the order ID as an int, handling string or number JSON types.
func (d *DomainDetails) OrderIDInt() (int, error) {
	return toInt(d.OrderID)
}

// Nameservers returns the non-empty nameservers as a slice.
func (d *DomainDetails) Nameservers() []string {
	fields := []string{d.NS1, d.NS2, d.NS3, d.NS4, d.NS5, d.NS6, d.NS7,
		d.NS8, d.NS9, d.NS10, d.NS11, d.NS12, d.NS13}
	var ns []string
	for _, v := range fields {
		if v != "" {
			ns = append(ns, v)
		}
	}
	return ns
}

// PrivacyEnabled returns whether privacy protection is active.
func (d *DomainDetails) PrivacyEnabled() bool {
	switch v := d.PrivacyProtected.(type) {
	case bool:
		return v
	case string:
		return strings.EqualFold(v, "true") || v == "1"
	case float64:
		return v != 0
	}
	return false
}

// ContactIDInt extracts an int from a contact ID field (which may be string or number).
func ContactIDInt(v interface{}) (int, error) {
	return toInt(v)
}

// ---------------------------------------------------------------------------
// Contact types
// ---------------------------------------------------------------------------

// ContactParams holds the fields for creating or updating a contact.
type ContactParams struct {
	CustomerID   int
	Type         string // "Contact", "UkContact", etc.
	Name         string
	Company      string
	Email        string
	AddressLine1 string
	AddressLine2 string
	AddressLine3 string
	City         string
	State        string
	Country      string
	Zipcode      string
	PhoneCC      string
	Phone        string
	FaxCC        string
	Fax          string
}

// ContactDetails holds the fields returned by the contacts/details endpoint.
type ContactDetails struct {
	ContactID     interface{} `json:"entityid"`
	Name          string      `json:"name"`
	Company       string      `json:"company"`
	Type          string      `json:"type"`
	Email         string      `json:"emailaddr"`
	PhoneCC       string      `json:"telnocc"`
	Phone         string      `json:"telno"`
	AddressLine1  string      `json:"address1"`
	AddressLine2  string      `json:"address2"`
	AddressLine3  string      `json:"address3"`
	City          string      `json:"city"`
	State         string      `json:"state"`
	Country       string      `json:"country"`
	Zipcode       string      `json:"zip"`
	CurrentStatus string      `json:"currentstatus"`
	CustomerID    interface{} `json:"customerid"`
}

// ContactIDInt returns the contact ID as an int.
func (c *ContactDetails) ContactIDInt() (int, error) {
	return toInt(c.ContactID)
}

// ---------------------------------------------------------------------------
// Domain search types
// ---------------------------------------------------------------------------

// DomainSearchResult holds a single domain from the search endpoint.
type DomainSearchResult struct {
	OrderID       interface{} `json:"orderid"`
	DomainName    string      `json:"domainname"`
	CurrentStatus string      `json:"currentstatus"`
	ProductKey    string      `json:"productkey"`
	ExpiryDate    string      `json:"endtime"`
	CreationDate  string      `json:"creationtime"`
	CustomerID    interface{} `json:"customerid"`
}

// DomainSearchResponse is the top-level response from search.json.
type DomainSearchResponse struct {
	RecordsCount int                  `json:"recsindb"`
	Results      []DomainSearchResult `json:"result"`
}

// DomainSearchParams holds optional filters for domain search.
type DomainSearchParams struct {
	NoOfRecords int
	PageNo      int
	Status      []string
	ProductKey  []string
	DomainName  string
}

// ---------------------------------------------------------------------------
// Availability types
// ---------------------------------------------------------------------------

// DomainAvailabilityResult holds the availability status for one domain.
type DomainAvailabilityResult struct {
	Domain string
	Status string // "available", "regthroughus", "regthroughothers", "unknown"
}

// ---------------------------------------------------------------------------
// Helper utilities
// ---------------------------------------------------------------------------

func toInt(v interface{}) (int, error) {
	switch val := v.(type) {
	case float64:
		return int(val), nil
	case string:
		return strconv.Atoi(val)
	case int:
		return val, nil
	case nil:
		return 0, fmt.Errorf("value is nil")
	default:
		return 0, fmt.Errorf("unexpected type %T", v)
	}
}

func (c *Client) authParams() url.Values {
	params := url.Values{}
	params.Set("auth-userid", strconv.Itoa(c.AuthUserID))
	params.Set("api-key", c.APIKey)
	return params
}

func (c *Client) getAndDecode(reqURL string, target interface{}) error {
	resp, err := c.HTTPClient.Get(reqURL)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}

	var errCheck struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}
	if err := json.Unmarshal(body, &errCheck); err != nil {
		return fmt.Errorf("parsing response: %w", err)
	}
	if errCheck.Status == "ERROR" {
		return fmt.Errorf("API error: %s", errCheck.Error)
	}

	return json.Unmarshal(body, target)
}

func (c *Client) postFormAndCheck(endpoint string, params url.Values) error {
	resp, err := c.HTTPClient.PostForm(c.BaseURL+endpoint, params)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}

	var result actionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("parsing response: %w", err)
	}
	if result.Status == "ERROR" {
		return fmt.Errorf("API error: %s", result.Error)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Domain — nameservers
// ---------------------------------------------------------------------------

// GetDomainNameservers fetches the current nameservers for a domain order.
func (c *Client) GetDomainNameservers(orderID int) ([]string, error) {
	params := c.authParams()
	params.Set("order-id", strconv.Itoa(orderID))
	params.Add("options", "NsDetails")

	var details DomainDetails
	if err := c.getAndDecode(fmt.Sprintf("%s/api/domains/details.json?%s", c.BaseURL, params.Encode()), &details); err != nil {
		return nil, err
	}
	return details.Nameservers(), nil
}

// ModifyNameservers updates the nameservers for a domain order.
func (c *Client) ModifyNameservers(orderID int, nameservers []string) error {
	params := c.authParams()
	params.Set("order-id", strconv.Itoa(orderID))
	for _, ns := range nameservers {
		params.Add("ns", ns)
	}
	return c.postFormAndCheck("/api/domains/modify-ns.json", params)
}

// ---------------------------------------------------------------------------
// Domain — details by name
// ---------------------------------------------------------------------------

// GetDomainByName fetches domain details (including order ID) by domain name.
func (c *Client) GetDomainByName(domainName string) (*DomainDetails, error) {
	params := c.authParams()
	params.Set("domain-name", domainName)
	params.Add("options", "OrderDetails")
	params.Add("options", "NsDetails")

	var details DomainDetails
	if err := c.getAndDecode(fmt.Sprintf("%s/api/domains/details-by-name.json?%s", c.BaseURL, params.Encode()), &details); err != nil {
		return nil, err
	}
	return &details, nil
}

// GetDomainDetails fetches full domain details by order ID with the given options.
func (c *Client) GetDomainDetails(orderID int, options []string) (*DomainDetails, error) {
	params := c.authParams()
	params.Set("order-id", strconv.Itoa(orderID))
	for _, o := range options {
		params.Add("options", o)
	}

	var details DomainDetails
	if err := c.getAndDecode(fmt.Sprintf("%s/api/domains/details.json?%s", c.BaseURL, params.Encode()), &details); err != nil {
		return nil, err
	}
	return &details, nil
}

// ---------------------------------------------------------------------------
// Domain — privacy protection
// ---------------------------------------------------------------------------

// ModifyPrivacyProtection enables or disables WHOIS privacy for a domain.
func (c *Client) ModifyPrivacyProtection(orderID int, protect bool, reason string) error {
	params := c.authParams()
	params.Set("order-id", strconv.Itoa(orderID))
	params.Set("protect-privacy", strconv.FormatBool(protect))
	params.Set("reason", reason)
	return c.postFormAndCheck("/api/domains/modify-privacy-protection.json", params)
}

// ---------------------------------------------------------------------------
// Domain — contacts
// ---------------------------------------------------------------------------

// ModifyDomainContacts assigns contact IDs to a domain registration order.
func (c *Client) ModifyDomainContacts(orderID, regContactID, adminContactID, techContactID, billingContactID int) error {
	params := c.authParams()
	params.Set("order-id", strconv.Itoa(orderID))
	params.Set("reg-contact-id", strconv.Itoa(regContactID))
	params.Set("admin-contact-id", strconv.Itoa(adminContactID))
	params.Set("tech-contact-id", strconv.Itoa(techContactID))
	params.Set("billing-contact-id", strconv.Itoa(billingContactID))
	return c.postFormAndCheck("/api/domains/modify-contact.json", params)
}

// GetDomainContactIDs returns the four contact IDs assigned to a domain.
func (c *Client) GetDomainContactIDs(orderID int) (reg, admin, tech, billing int, err error) {
	details, err := c.GetDomainDetails(orderID, []string{"ContactIds"})
	if err != nil {
		return
	}
	reg, err = ContactIDInt(details.RegistrantContactID)
	if err != nil {
		return
	}
	admin, err = ContactIDInt(details.AdminContactID)
	if err != nil {
		return
	}
	tech, err = ContactIDInt(details.TechContactID)
	if err != nil {
		return
	}
	billing, err = ContactIDInt(details.BillingContactID)
	return
}

// ---------------------------------------------------------------------------
// Domain — child nameservers (glue records)
// ---------------------------------------------------------------------------

// AddChildNameserver creates a child nameserver (glue record) under a domain.
func (c *Client) AddChildNameserver(orderID int, hostname string, ips []string) error {
	params := c.authParams()
	params.Set("order-id", strconv.Itoa(orderID))
	params.Set("cns", hostname)
	for _, ip := range ips {
		params.Add("ip", ip)
	}
	return c.postFormAndCheck("/api/domains/add-cns.json", params)
}

// ModifyChildNameserverHostname renames a child nameserver.
func (c *Client) ModifyChildNameserverHostname(orderID int, oldHostname, newHostname string) error {
	params := c.authParams()
	params.Set("order-id", strconv.Itoa(orderID))
	params.Set("old-cns", oldHostname)
	params.Set("new-cns", newHostname)
	return c.postFormAndCheck("/api/domains/modify-cns-name.json", params)
}

// ModifyChildNameserverIP changes one IP address of a child nameserver.
func (c *Client) ModifyChildNameserverIP(orderID int, hostname, oldIP, newIP string) error {
	params := c.authParams()
	params.Set("order-id", strconv.Itoa(orderID))
	params.Set("cns", hostname)
	params.Set("old-ip", oldIP)
	params.Set("new-ip", newIP)
	return c.postFormAndCheck("/api/domains/modify-cns-ip.json", params)
}

// DeleteChildNameserver removes a child nameserver from a domain.
func (c *Client) DeleteChildNameserver(orderID int, hostname string) error {
	params := c.authParams()
	params.Set("order-id", strconv.Itoa(orderID))
	params.Set("cns", hostname)
	return c.postFormAndCheck("/api/domains/delete-cns.json", params)
}

// GetChildNameservers returns a map of hostname -> []IPs for a domain's child NS records.
func (c *Client) GetChildNameservers(orderID int) (map[string][]string, error) {
	details, err := c.GetDomainDetails(orderID, []string{"NsDetails"})
	if err != nil {
		return nil, err
	}

	result := map[string][]string{}
	for host, raw := range details.CNS {
		switch v := raw.(type) {
		case string:
			result[host] = []string{v}
		case []interface{}:
			var ips []string
			for _, ip := range v {
				if s, ok := ip.(string); ok {
					ips = append(ips, s)
				}
			}
			result[host] = ips
		}
	}
	return result, nil
}

// ---------------------------------------------------------------------------
// Contacts
// ---------------------------------------------------------------------------

func contactParams(p ContactParams) url.Values {
	params := url.Values{}
	params.Set("name", p.Name)
	params.Set("company", p.Company)
	params.Set("email", p.Email)
	params.Set("address-line-1", p.AddressLine1)
	if p.AddressLine2 != "" {
		params.Set("address-line-2", p.AddressLine2)
	}
	if p.AddressLine3 != "" {
		params.Set("address-line-3", p.AddressLine3)
	}
	params.Set("city", p.City)
	if p.State != "" {
		params.Set("state", p.State)
	}
	params.Set("country", p.Country)
	params.Set("zipcode", p.Zipcode)
	params.Set("phone-cc", p.PhoneCC)
	params.Set("phone", p.Phone)
	if p.FaxCC != "" {
		params.Set("fax-cc", p.FaxCC)
	}
	if p.Fax != "" {
		params.Set("fax", p.Fax)
	}
	return params
}

// AddContact creates a new contact and returns its contact ID.
func (c *Client) AddContact(p ContactParams) (int, error) {
	params := c.authParams()
	for k, v := range contactParams(p) {
		params[k] = v
	}
	params.Set("customer-id", strconv.Itoa(p.CustomerID))
	params.Set("type", p.Type)

	resp, err := c.HTTPClient.PostForm(c.BaseURL+"/api/contacts/add.json", params)
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("reading response: %w", err)
	}

	// Success response is a plain integer (the new contact-id).
	// Error response is a JSON object with status=ERROR.
	var errCheck struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}
	if jsonErr := json.Unmarshal(body, &errCheck); jsonErr == nil && errCheck.Status == "ERROR" {
		return 0, fmt.Errorf("API error: %s", errCheck.Error)
	}

	id, err := strconv.Atoi(strings.TrimSpace(string(body)))
	if err != nil {
		return 0, fmt.Errorf("unexpected response (expected contact ID integer): %s", string(body))
	}
	return id, nil
}

// ModifyContact updates an existing contact's details.
func (c *Client) ModifyContact(contactID int, p ContactParams) error {
	params := c.authParams()
	for k, v := range contactParams(p) {
		params[k] = v
	}
	params.Set("contact-id", strconv.Itoa(contactID))
	return c.postFormAndCheck("/api/contacts/modify.json", params)
}

// GetContact fetches a contact's details by ID.
func (c *Client) GetContact(contactID int) (*ContactDetails, error) {
	params := c.authParams()
	params.Set("contact-id", strconv.Itoa(contactID))

	var details ContactDetails
	if err := c.getAndDecode(fmt.Sprintf("%s/api/contacts/details.json?%s", c.BaseURL, params.Encode()), &details); err != nil {
		return nil, err
	}
	return &details, nil
}

// DeleteContact removes a contact.
func (c *Client) DeleteContact(contactID int) error {
	params := c.authParams()
	params.Set("contact-id", strconv.Itoa(contactID))
	return c.postFormAndCheck("/api/contacts/delete.json", params)
}

// SearchContacts searches contacts for a customer, returning up to noOfRecords on pageNo.
func (c *Client) SearchContacts(customerID, noOfRecords, pageNo int, name, email, contactType string) ([]ContactDetails, int, error) {
	params := c.authParams()
	params.Set("customer-id", strconv.Itoa(customerID))
	params.Set("no-of-records", strconv.Itoa(noOfRecords))
	params.Set("page-no", strconv.Itoa(pageNo))
	if name != "" {
		params.Set("name", name)
	}
	if email != "" {
		params.Set("email", email)
	}
	if contactType != "" {
		params.Set("type", contactType)
	}

	var raw map[string]json.RawMessage
	if err := c.getAndDecode(fmt.Sprintf("%s/api/contacts/search.json?%s", c.BaseURL, params.Encode()), &raw); err != nil {
		return nil, 0, err
	}

	total := 0
	if v, ok := raw["recsindb"]; ok {
		_ = json.Unmarshal(v, &total)
	}

	var contacts []ContactDetails
	if v, ok := raw["result"]; ok {
		_ = json.Unmarshal(v, &contacts)
	}
	return contacts, total, nil
}

// ---------------------------------------------------------------------------
// Domain search
// ---------------------------------------------------------------------------

// SearchDomains returns domain orders matching the given parameters.
func (c *Client) SearchDomains(p DomainSearchParams) ([]DomainSearchResult, int, error) {
	params := c.authParams()
	if p.NoOfRecords <= 0 {
		p.NoOfRecords = 50
	}
	if p.PageNo <= 0 {
		p.PageNo = 1
	}
	params.Set("no-of-records", strconv.Itoa(p.NoOfRecords))
	params.Set("page-no", strconv.Itoa(p.PageNo))
	for _, s := range p.Status {
		params.Add("status", s)
	}
	for _, k := range p.ProductKey {
		params.Add("product-key", k)
	}
	if p.DomainName != "" {
		params.Set("domain-name", p.DomainName)
	}

	var raw map[string]json.RawMessage
	if err := c.getAndDecode(fmt.Sprintf("%s/api/domains/search.json?%s", c.BaseURL, params.Encode()), &raw); err != nil {
		return nil, 0, err
	}

	total := 0
	if v, ok := raw["recsindb"]; ok {
		_ = json.Unmarshal(v, &total)
	}

	var results []DomainSearchResult
	if v, ok := raw["result"]; ok {
		_ = json.Unmarshal(v, &results)
	}
	return results, total, nil
}

// ---------------------------------------------------------------------------
// Domain availability
// ---------------------------------------------------------------------------

// CheckDomainAvailability checks whether domain+TLD combinations are available.
// Returns a map of "domain.tld" -> status string.
func (c *Client) CheckDomainAvailability(domains, tlds []string) (map[string]string, error) {
	params := c.authParams()
	for _, d := range domains {
		params.Add("domain-name", d)
	}
	for _, t := range tlds {
		params.Add("tlds", t)
	}

	// The availability check endpoint uses a different subdomain.
	checkURL := strings.Replace(c.BaseURL, "api.", "domaincheck.", 1)
	if checkURL == c.BaseURL {
		// Fallback: try inserting domaincheck subdomain if no "api." prefix found.
		checkURL = strings.Replace(c.BaseURL, "://", "://domaincheck.", 1)
	}

	var raw map[string]json.RawMessage
	if err := c.getAndDecode(fmt.Sprintf("%s/api/domains/available.json?%s", checkURL, params.Encode()), &raw); err != nil {
		return nil, err
	}

	result := map[string]string{}
	for key, val := range raw {
		var entry struct {
			Status string `json:"status"`
		}
		if err := json.Unmarshal(val, &entry); err == nil {
			result[key] = entry.Status
		}
	}
	return result, nil
}

// ---------------------------------------------------------------------------
// Default nameservers
// ---------------------------------------------------------------------------

// GetDefaultNameservers returns the customer's configured default nameservers.
func (c *Client) GetDefaultNameservers(customerID int) ([]string, error) {
	params := c.authParams()
	params.Set("customer-id", strconv.Itoa(customerID))

	var ns []string
	if err := c.getAndDecode(fmt.Sprintf("%s/api/domains/customer-default-ns.json?%s", c.BaseURL, params.Encode()), &ns); err != nil {
		return nil, err
	}
	return ns, nil
}
