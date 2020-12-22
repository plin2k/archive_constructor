# archive_constructor
 Archive Constructor from JSON


go run archive_constructor.go '<json>'



```json
{
   "archive_name":"name",
   "output_path":"foo/bar",
   "data":{
      "folders":[
         {
            "name":"folder_name_1",
            "folders":[
                {
                   "name":"folder_name_1",
                   "folders":[
                      
                   ],
                   "sources":[
                      {
                         "name":"file_name_1",
                         "url":"https://plin2k.org",
                         "extention":"pdf"
                      },
                      {
                         "name":"file_name_2",
                         "url":"https://plin2k.org",
                         "extention":"pdf"
                      }
                   ]
                }
            ],
            "sources":[
               {
                  "name":"file_name_1",
                  "url":"https://plin2k.org",
                  "extention":"pdf"
               },
               {
                  "name":"file_name_2",
                  "url":"https://plin2k.org",
                  "extention":"pdf"
               }
            ]
         },
         {
            "name":"folder_name_2",
            "folders":[
               
            ],
            "sources":[
               {
                  "name":"file_name_3",
                  "url":"https://plin2k.org",
                  "extention":"pdf"
               },
               {
                  "name":"file_name_4",
                  "url":"https://plin2k.org",
                  "extention":"pdf"
               }
            ]
         }
      ],
      "sources":[
         
      ]
   }
}
```
