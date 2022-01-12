# go gql-server


## テンプレートコード生成 (初回)


graph/schema.graphqlsを記述する

```bash
$ ~/go/bin/gqlgen init
$ go mod tidy
```

```
.
├── README.md
├── docker-compose.yaml
├── go.mod
├── go.sum
├── gqlgen.yml
├── graph
│   ├── generated
│   │   └── generated.go
│   ├── model
│   │   └── models_gen.go
│   ├── resolver.go
│   ├── schema.graphqls
│   └── schema.resolvers.go
├── postgres
│   ├── Dockerfile
│   └── init
│       ├── 1_create.sql
│       ├── 2_copy.sql
│       ├── company20200619.csv
│       ├── join20211208.csv
│       ├── line20211208free.csv
│       └── station20211222free.csv
└── server.go
```

## テンプレートコード生成 (2回目以降)

SQLの結合が必要になるフィールドにはGraphQL resolverを分離し別に用意する。

(つまり、gqlgen側にクエリの実行順序やレスポンス整形を委ねる)

これを行うと各resolverの実装をシンプルに保ちつつ、複雑なクエリに対応できる。


```yaml
models:
  # 中略
  Station:                # GraphQLスキーマのStation型
    fields:
      beforeStation:      # フィールド名
        resolver: true
      afterStation:       # フィールド名
        resolver: true
      transferStation:    # フィールド名
        resolver:  true
```

```bash
$ ~/go/bin/gqlgen generate
```



## DBスキーマからstructを生成


Generate code from a custom SQL query for Postgres


```bash
$ mkdir -p models
$ # 駅CD検索
$ ~/go/bin/xo query 'postgresql://root:postgres@127.0.0.1:5433/root?sslmode=disable' -M -B -T StationByCD -o models/ << ENDSQL
select l.line_cd, l.line_name, s.station_cd, station_g_cd, s.station_name, s.address
from station s
         inner join line l on s.line_cd = l.line_cd
where s.station_cd = %%stationCD int%%
  and s.e_status = 0
ENDSQL
$ # 駅名検索
$ ~/go/bin/xo query 'postgresql://root:postgres@127.0.0.1:5433/root?sslmode=disable' -M -B -T StationByName -o models/ << ENDSQL
select l.line_cd, l.line_name, s.station_cd, station_g_cd, s.station_name, s.address
from station s
         inner join line l on s.line_cd = l.line_cd
where s.station_name = %%stationName string%%
  and s.e_status = 0
ENDSQL
$ # 駅検索
$ ~/go/bin/xo query 'postgresql://root:postgres@127.0.0.1:5433/root?sslmode=disable' -M -B -T Station -o models/ << ENDSQL
select l.line_cd, l.line_name, s.station_cd, station_g_cd, s.station_name, s.address
from station s
         inner join line l on s.line_cd = l.line_cd
where s.station_cd = %%stationCD int%%
  and s.e_status = 0
ENDSQL
$ # 隣駅(after)検索
$ ~/go/bin/xo query 'postgresql://root:postgres@127.0.0.1:5433/root?sslmode=disable' -M -B -T After -o models/ << ENDSQL
select sl.line_cd,
       sl.line_name,
       s.station_cd,
       s.station_name,
       s.address,
       COALESCE(js.station_cd, 0)    as after_station_cd,
       COALESCE(js.station_name, '') as after_station_name,
       COALESCE(js.station_g_cd, 0)  as after_station_g_cd,
       COALESCE(js.address, '')      as after_station_address
from station s
         left outer join line sl on s.line_cd = sl.line_cd
         left outer join station_join j on s.line_cd = j.line_cd and s.station_cd = j.station_cd2
         left outer join station js on j.station_cd1 = js.station_cd
where s.e_status = 0
  and s.station_cd = %%stationCD int%%
ENDSQL
$ # 隣駅(before)検索
$ ~/go/bin/xo query 'postgresql://root:postgres@127.0.0.1:5433/root?sslmode=disable' -M -B -T Before -o models/ << ENDSQL
select sl.line_cd,
       sl.line_name,
       s.station_cd,
       s.station_name,
       s.address,
       COALESCE(js.station_cd, 0)    as before_station_cd,
       COALESCE(js.station_name, '') as before_station_name,
       COALESCE(js.station_g_cd, 0)  as before_station_g_cd,
       COALESCE(js.address, '')      as before_station_address
from station s
         left outer join line sl on s.line_cd = sl.line_cd
         left outer join station_join j on s.line_cd = j.line_cd and s.station_cd = j.station_cd1
         left outer join station js on j.station_cd2 = js.station_cd
where s.e_status = 0
  and s.station_cd = %%stationCD int%%
ENDSQL
$ # 乗り換え駅検索
$ ~/go/bin/xo query 'postgresql://root:postgres@127.0.0.1:5433/root?sslmode=disable' -M -B -T Transfer -o models/ << ENDSQL
select s.station_cd,
       ls.line_cd,
       ls.line_name,
       s.station_name,
       s.station_g_cd,
       s.address,
       COALESCE(lt.line_cd, 0)     as transfer_line_cd,
       COALESCE(lt.line_name, '')   as transfer_line_name,
       COALESCE(t.station_cd, 0)   as transfer_station_cd,
       COALESCE(t.station_name, '') as transfer_station_name,
       COALESCE(t.address, '')      as transfer_address
from station s
         left outer join station t on s.station_g_cd = t.station_g_cd and s.station_cd <> t.station_cd
         left outer join line ls on s.line_cd = ls.line_cd
         left outer join line lt on t.line_cd = lt.line_cd
where s.station_cd = %%stationCD int%%
ENDSQL
```


## sqlを送る

```graphql
query osaki{
  stationByName(stationName: "大崎") {
    lineName
    stationCD
    stationName
  }
}
```

```graphql
query nextStation {
  stationByCD(stationCD: 1130201) {
    lineName
    stationCD
    stationName
    beforeStation {
      lineName
      stationCD
      stationName
    }
    afterStation {
      lineName
      stationCD
      stationName
    }
  }
}
```

```graphql
query stationByCD {
  stationByCD(stationCD: 1130201) {
    lineName
    stationCD
    stationName
    beforeStation {
      lineName
      stationCD
      stationName
      transferStation {
        lineName
        stationCD
        stationName
      }
    }
    afterStation {
      lineName
      stationCD
      stationName
    }
  }
}
```

```graphql
fragment stationF on Station {
  lineName
  stationCD
  stationName
}

query stationByCD {
  stationByCD(stationCD: 1130201) {
    ...stationF
    beforeStation {
      ...stationF
      transferStation {
        ...stationF
      }
    }
    afterStation {
      ...stationF
    }
  }
}
```