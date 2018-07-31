# git init


Simple Amo CRM API wrapper.

## Features


* Get Account
* Add Companies
* Add Contacts
* Add Tasks
* Add Notes  

In future API can be changed/modified. 
Feel free to fork and modify for your needs. 

## Example

```go

client, err := NewClient("https://yourdomain.amocrm.ru", nil)
if err != nil {
    panic(err)
}

auth, err := client.Authorize("your@email.com", "amo_api_key")
if err != nil {
    panic(err)
}

if !auth.Response.Auth {
    panic("NO AUTH")
}

acc, err := client.GetAccount(false, []string{"users", "custom_fields"})
if err != nil {
    panic(err)
}

compId, err := client.AddCompany(Company{
    Name:              "Company Name",
    ResponsibleUserID: 123456,
    CreatedBy:         123456,
    // Information about your custom fields and IDs 
    // you can get from account
    CustomFields: []CustomField{
        {
            ID: 317255, // Phone
            Values: []CSValue{
                {Value: MultiValuesString("111-222-333"), Enum: "702113"},
            },
        },
    },
})
if err != nil {
    panic(err)
}

```