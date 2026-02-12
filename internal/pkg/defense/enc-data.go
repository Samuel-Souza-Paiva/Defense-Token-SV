package defense

import "encoding/json"

type EncData struct{
  Payload EncDataPayload
  Res EncDataRes
}

func (encData *EncData) CreatePayload(username string) {
  encData.Payload = EncDataPayload{
    UserName: username,
    ClientType: "WINPC_V2",
  }
}

func (encData *EncData) SetRes(data []byte) {
  var res EncDataRes

  err := json.Unmarshal(data, &res)

  if err != nil {
    return
  }

  encData.Res = res
}
