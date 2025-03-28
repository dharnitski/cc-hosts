# Data folder contain Common Crawl Hosts Web Graph

## How to get the data

Login into into AWS and navigate to s3://commoncrawl/projects/hyperlinkgraph/ folder. Find snapshot you want to use. In my case it is s3://commoncrawl/projects/hyperlinkgraph/cc-main-2024-oct-nov-dec/. Note: the latest data may be not in the latest folder because name includes month name instead of month number.

Download and uncompress `.txt.gz` files from S3 `s3://commoncrawl/projects/hyperlinkgraph/cc-main-XXX/host/vertices/` folder into  `data/vertices`. 

Download and uncompress `.txt.gz` files from S3 `s3://commoncrawl/projects/hyperlinkgraph/cc-main-XXX/host/edges/` folder into  `data/edges`.

Finals result

```
./data/vertices
./data/vertices/part-00004-4ba7987d-67a0-4f7d-b410-1d92df440699-c000.txt
./data/vertices/part-00005-4ba7987d-67a0-4f7d-b410-1d92df440699-c000.txt
./data/vertices/part-00000-4ba7987d-67a0-4f7d-b410-1d92df440699-c000.txt
./data/vertices/part-00001-4ba7987d-67a0-4f7d-b410-1d92df440699-c000.txt
./data/vertices/part-00002-4ba7987d-67a0-4f7d-b410-1d92df440699-c000.txt
./data/vertices/part-00003-4ba7987d-67a0-4f7d-b410-1d92df440699-c000.txt
./data/edges
./data/edges/part-00011-02106921-c60f-49b6-912c-b03ea5690455-c000.txt
./data/edges/part-00006-02106921-c60f-49b6-912c-b03ea5690455-c000.txt
./data/edges/part-00010-02106921-c60f-49b6-912c-b03ea5690455-c000.txt
./data/edges/part-00007-02106921-c60f-49b6-912c-b03ea5690455-c000.txt
./data/edges/part-00008-02106921-c60f-49b6-912c-b03ea5690455-c000.txt
./data/edges/part-00004-02106921-c60f-49b6-912c-b03ea5690455-c000.txt
./data/edges/part-00009-02106921-c60f-49b6-912c-b03ea5690455-c000.txt
./data/edges/part-00005-02106921-c60f-49b6-912c-b03ea5690455-c000.txt
./data/edges/part-00002-02106921-c60f-49b6-912c-b03ea5690455-c000.txt
./data/edges/part-00003-02106921-c60f-49b6-912c-b03ea5690455-c000.txt
./data/edges/part-00000-02106921-c60f-49b6-912c-b03ea5690455-c000.txt
./data/edges/part-00001-02106921-c60f-49b6-912c-b03ea5690455-c000.txt
```

Files structure

Vertices are stored as two columns delimited by tab. 
First column is VerticeID. It is unique incremented number.
Second column is DNS eTLD+1 name in reversed DNS format.
Vertices are sorted by ID using int (not string) comparison logic.  
Every file is continuation of data and IDs are incremented from file to file.


```
0	aaa.11111
1	aaa.3
2	aaa.a
3	aaa.aa
4	aaa.aa.aaaa
5	aaa.aaa
6	aaa.aaa.242.133
7	aaa.aaa.aaa
8	aaa.aaa.aaa.aaa
9	aaa.aaa.fsdfsd
10	aaa.aaaa
11	aaa.aaaa.aaa
12	aaa.aaaa.aaaaaaaaa
13	aaa.aaaa.aaaaaaaaa.aaa
14	aaa.aaaaa
15	aaa.aaaaaa
16	aaa.aaaaaaa.aaa
17	aaa.aaaaaaaa.aaa
18	aaa.aaaaaaaaa
19	aaa.aaaaaaaaa.aaa
20	aaa.aaaaaaaaaaaaaaaaaaaa
21	aaa.aaaaaaaaaaaaaaaaaaaaaaaaaaa.aaa
22	aaa.acg.locator
23	aaa.compromises

```

Edges are stored as two Vertice ID values delimited by tab.
First Vertice identifies host that contains the link.
Second Vertice is host pointed by that link.
Edges are sorted by first Vertice ID using int (not string) comparison logic.
Each file is not continuation of index. Instead. Each file contains subset of all Vertices. 


```
75	63216723
75	229821733
77	47814421
84	40218536
84	40219361
84	119920069
84	205993715
84	277542382
90	40219361
90	219559011
90	229862070
96	47814421
100	40219273
111	138
111	32849676
111	91890049
112	138
112	32849676
112	92673207
114	138
114	3323038
114	32849676
114	91890049
```


## Create Reversed 

Reverse Vertices in Edges file and save into `data/edges_reversed` folder.

```
$python reverse.py  
```

Data is not sorted after Vertices are reversed.
Run script to sort files. Sorting can take several hours.

```
$./sort.sh  
``

Sorted Reversed Edges data 

```
17	90104917
34	124288668
34	256919968
51	212011488
63	252604196
69	161511398
69	188287621
69	189122683
69	192815673
69	243778869
69	40219320
69	45450131
69	92659650
71	133263766
71	165548565
71	40219320
71	40219391
71	66062753
73	256919968
```


## Scripts
