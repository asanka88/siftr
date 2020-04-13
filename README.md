[![MIT licensed](https://img.shields.io/badge/license-MIT-green.svg)](./LICENSE)

# Siftr

Credit for expand and flatten goes to [hashicorp/terraform](https://github.com/hashicorp/terraform/tree/master/flatmap). Just played around with that logic a bit to enhance the performance some.

### Supported actions
#### Whitelist
Check the given JSON payload against a list of allowable fields and remove those that do not belong.

#### Sibling
Check the value of the given field against the list of allowable values and remove the entire object if it does not belong.
So given this policy:
```json
{
  "sibling": {
      "woods.hardwoods.[*].commonName": [
        "OAK",
        "MP"
      ]
  }
}
```
And this data:
```json
{
  "woods": {
    "hardwoods": [
      {
        "id": "oak",
        "commonName": "OAK",
        "genus": "quercus",
        "types": [
          {
            "id": "redOak",
            "commonName": "Red Oak",
            "hardness": 1290
          }
        ]
      },
      {
        "id": "walnut",
        "commonName": "WN",
        "genus": "juglans",
        "types": [
          {
            "id": "black",
            "commonName": "Black Walnut",
            "hardness": 1010
          }
        ]
      }
    ]
  }
}
```
Would result in (with the walnut - WN object removed):
```json
{
  "woods": {
    "hardwoods": [
      {
        "id": "oak",
        "commonName": "OAK",
        "genus": "quercus",
        "types": [
          {
            "id": "redOak",
            "commonName": "Red Oak",
            "hardness": 1290
          }
        ]
      }
    ]
  }
}
```