package eve

type Killmail struct {
	Package struct {
		KillID int
		Killmail struct {
			Victim struct {
				CorporationID int `json:"corporation_id"`
			}
			Attackers []struct {
				CorporationID int `json:"corporation_id"`
			}
		}
	}
}
