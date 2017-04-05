package peergen

import (
	"encoding/json"
	"fmt"
	"github.com/ipcjk/ixgen/ixtypes"
	"io"
	"log"
	"strconv"
	"time"
)

// Functions that convert our Ix to different JSON / XML - configurations
// ConvertIxToJuniperJSON -> Juniper JSON
// ConvertIxToBrocadeSlxJSON -> Brocade SLX

func (p *Peergen) ConvertIxToJson(ix ixtypes.Ix, w io.Writer) {
	res, err := json.Marshal(ix)
	if err != nil {
		log.Fatalf("Cant decode IX into native format: err")
	}
	fmt.Fprintf(w, string(res))
}

func (p *Peergen) ConvertIxToJsonPretty(ix ixtypes.Ix, w io.Writer) {
	res, err := json.MarshalIndent(ix, "", "\t")
	if err != nil {
		log.Fatalf("Cant decode IX into native format: err")
	}
	fmt.Fprintf(w, string(res))
}

func (p *Peergen) ConvertIxToBrocadeSlxJSON(ix ixtypes.Ix, w io.Writer) {
	log.Fatal("Not done yet")
}

func (p *Peergen) ConvertIxToJuniperJSON(ix ixtypes.Ix, w io.Writer) {
	var junosConfiguration = junOsJSON{
		[]junosConfiguration{
			{
				junosAttributes{},
				[]junosBGPProtocol{
					{
						[]junosBgpGroup{
							{[]junosGroup{}},
						},
					},
				},
			},
		},
	}

	for k := range ix.PeeringGroups {
		var junosPeerConfiguration junosGroup
		for i := range ix.PeersReady {
			if ix.PeersReady[i].Group == k {
				junosPeerConfiguration.Name = junosDataString{k}
				junosPeerConfiguration.Type = []struct{ junosDataString }{
					{junosDataString{"external"}},
				}
				if ix.PeersReady[i].Ipv4Enabled {
					junosPeerConfiguration.Neighbor = append(junosPeerConfiguration.Neighbor,
						junosNeighbor{
							Family: junosFamily{
								{Inet6: []junosFamilyInet6{},
									Inet: []junosFamilyInet4{
										{InetUnicast: []junosLabeledUnicast{
											{
												PrefixLimit: []junosPrefixLimit{
													{Maximum: []junosMaximumLimit{
														{junosDataInt64String: junosDataInt64String{Data: strconv.FormatInt(ix.PeersReady[i].InfoPrefixes4, 10)}},
													}},
												},
											},
										}},
									}},
							},
							Name:   junosDataIP{Data: ix.PeersReady[i].Ipv4Addr},
							PeerAs: []junosDataInt64String{{Data: ix.PeersReady[i].ASN}},
						})
				}
				if ix.PeersReady[i].Ipv6Enabled {
					junosPeerConfiguration.Neighbor = append(junosPeerConfiguration.Neighbor,
						junosNeighbor{
							Family: junosFamily{
								{Inet6: []junosFamilyInet6{{Inet6Unicast: []junosLabeledUnicast{
									{
										PrefixLimit: []junosPrefixLimit{
											{Maximum: []junosMaximumLimit{
												{junosDataInt64String: junosDataInt64String{Data: strconv.FormatInt(ix.PeersReady[i].InfoPrefixes6, 10)}},
											}},
										},
									},
								}}},
									Inet: []junosFamilyInet4{}},
							},
							Name:   junosDataIP{Data: ix.PeersReady[i].Ipv6Addr},
							PeerAs: []junosDataInt64String{{Data: ix.PeersReady[i].ASN}},
						})
				}
				if len(junosPeerConfiguration.Neighbor) > 0 {
					junosConfiguration.Configuration[0].Protocols[0].Bgp[0].Group = append(
						junosConfiguration.Configuration[0].Protocols[0].Bgp[0].Group,
						junosPeerConfiguration)
				}
			}
		}
	}

	/* Add timestamp */
	t := time.Now()
	junosConfiguration.Configuration[0].Attributes.JunosChangedSeconds = t.Unix()
	junosConfiguration.Configuration[0].Attributes.JunosChangedLocaltime = t.String()

	res, err := json.Marshal(junosConfiguration)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, string(res))
}

/*
var junosConfiguration = junOsJSON{
	[]junosConfiguration{
		{
			junosAttributes{},
			[]junosBGPProtocol{
				{
					[]junosBgpGroup{
						{[]junosGroup{
							junosGroup{
								junosDataString{"Group"},
								[]junosNeighbor{
									{junosDataIP{},
									 []junosDataInt64String{},
									},

								},
							},
						}},
					},
				},
			},
		},
	},
}

	/* Generate Group for peers without Group */
/*
var wildgroup = junosGroup{
	junosDataString{"Group"},
	[]junosNeighbor{
		{junosDataIP{net.IP("127.0.0.1")},
			[]junosDataInt64String{},
		},

	},
}
*/
