package types

import "github.com/near/borsh-go"

type Data struct {
	Name                 string      `json:"name"`
	Symbol               string      `json:"symbol"`
	Uri                  string      `json:"uri"`
	SellerFeeBasisPoints uint16      `json:"seller_fee_basis_points,omitempty"`
	Creators             *[]Creator  `json:"creators"`
	Collection           *Collection `json:"collection"`
	Uses                 *Uses
}

type Metadata struct {
	Name                 string      `json:"name"`
	Symbol               string      `json:"symbol"`
	Description          string      `json:"description"`
	Image                string      `json:"image"`
	SellerFeeBasisPoints int         `json:"seller_fee_basis_points,omitempty"`
	Properties           Properties  `json:"properties"`
	Collection           Collection  `json:"collection"`
	Attributes           []Attribute `json:"attributes"`
}

type Properties struct {
	Files    []Files   `json:"files"`
	Creators []Creator `json:"creators"`
}

type Files struct {
	Uri  string `json:"uri"`
	Type string `json:"type"`
}

type Creator struct {
	Address string `json:"address"`
	Share   uint8  `json:"share"`
}

type Collection struct {
	Name   string `json:"name"`
	Family string `json:"family"`
}

type Attribute struct {
	Type  string `json:"trait_type"`
	Value string `json:"value"`
}

type Uses struct {
	UseMethod borsh.Enum
	Remaining uint64
	Total     uint64
}
