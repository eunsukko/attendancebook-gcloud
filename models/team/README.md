# team_generation.json
팀을 쉽게 만들기 위함

[
    {  
      "name":"은석이네",
      "member_infos":[  
        {  
          "name":"고은석",
          "birth":"920727",
          "tags":[  
            "고씨",
            "프로그래머"
          ]
        },
        {  
          "name":"고경석"
        }
      ],
      "member_valid_period":{  
        "from":"2018-11-04T00:00:00Z",
        "to":"2019-11-04T00:00:00Z"
      }
    },
    {  
      "name":"2019년 믿음교구 1셀",
      "member_infos":[  
        {  
          "name":"이사라",
          "birth":"950920",
          "tags":[  
            "셀리더"
          ]
        },
        {  
          "name":"김민수",
          "birth":"940913",
          "tags":[  
            "부셀리더"
          ]
        },
        {  
          "name":"김윤호",
          "birth":"990428"
        },
        {  
          "name":"정윤선",
          "birth":"970907"
        },
        {  
          "name":"천도현",
          "birth":"950813"
        },
        {  
          "name":"조병민",
          "birth":"990611"
        },
        {  
          "name":"고은석",
          "birth":"920727"
        },
        {  
          "name":"이재성",
          "birth":"971009"
        }
      ],
      "member_valid_period":{  
        "from":"2018-11-04T00:00:00Z",
        "to":"2019-11-04T00:00:00Z"
      }
    },
]


## 주의사항
member 중에 동일한 이름을 가진 사람이 또 있을 경우를 위해서 name + birth 를 제공
name 만 넣으면 해당 name 을 가진 모든 사람을 추가하고
name + birth 를 넣으면 해당 정보를 모두 가지는 사람을 추가하도록
(지금은 name + birth 로 유일한 사람을 표현한다고 가정하기에 tags 같은 정보를 함께 제공, 즉 그 사람을 온전히 찾는다는 가정이 있는 것임)


