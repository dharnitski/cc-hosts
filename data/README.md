# Data folder containing Common Crawl Hosts Web Graph


Data is collected by commoncrawl.org and hosted in public `commoncrawl` s3 bucket.
Anonymous access is disabled to that bucket, but any authenticated user can download data.

Check this link for more details - https://commoncrawl.org/get-started

## Download Common Crawl files

### Find latests data identifier

Common Crawls constantly scrubs internet and regularly publishes datasets with collected data.
For Web Graphs that currently happens every 3 month. 

We need to find the latest dataset to use. 

Go to https://index.commoncrawl.org/graphinfo.json and find latest `id` and `location`.

At the this documentation is written it is `cc-main-2025-feb-mar-apr` and `s3://commoncrawl/projects/hyperlinkgraph/cc-main-2025-feb-mar-apr/`

Alternative path to get that data:

Go to https://commoncrawl.org/web-graphs and pick it in that page. You will be redirected to page that contains downloads for that dataset.
In my case it is https://data.commoncrawl.org/projects/hyperlinkgraph/cc-main-2025-feb-mar-apr/index.html.  `cc-main-2025-feb-mar-apr` is what we need. That is dataset identifier. It is combined dataset containing data from Febriary, March and April of 2025. That page also has direct links to download data data in  [BVGraph](https://webgraph.di.unimi.it/docs/it/unimi/dsi/webgraph/BVGraph.html) format. 


### Download data

Hosts Web Graph data stored in S3 folder `s3://commoncrawl/projects/hyperlinkgraph/<Release ID>/host`.

Commands to download edges and vertices:

```bash
aws s3 cp s3://commoncrawl/projects/hyperlinkgraph/<Release ID>/host/edges/ data/edges --recursive

aws s3 cp s3://commoncrawl/projects/hyperlinkgraph/<Release ID>/host/vertices/ data/vertices --recursive
```

For release `cc-main-2025-feb-mar-apr` they looks like this:

```bash
aws s3 cp s3://commoncrawl/projects/hyperlinkgraph/cc-main-2025-jan-feb-mar/host/edges/ data/edges --recursive

aws s3 cp s3://commoncrawl/projects/hyperlinkgraph/cc-main-2025-jan-feb-mar/host/vertices/ data/vertices --recursive
```

These commands downloads about 14 GB for edges and 3 GB for vertices of `.txt.gz` files. 


## Extract files from archives

`gz` saves lots of space but we cannot use it directly. That format does not allow reading data by index.

Use this command to extract data:

```bash
find data -type f -name "*.txt.gz" -exec gunzip {} \;
```

Now we have about 55 GB for edges and 10 GB for vertices in `.txt` files.

## Finals result

```
./data/edges:
part-00000-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00001-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00002-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00003-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00004-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00005-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00006-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00007-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00008-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00009-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00010-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00011-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00012-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00013-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00014-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00015-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00016-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00017-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00018-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00019-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00020-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00021-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00022-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00023-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt

./data/vertices:
part-00000-b2128a57-30c4-4114-9ac0-344b81f88dfe-c000.txt
part-00001-b2128a57-30c4-4114-9ac0-344b81f88dfe-c000.txt
part-00002-b2128a57-30c4-4114-9ac0-344b81f88dfe-c000.txt
part-00003-b2128a57-30c4-4114-9ac0-344b81f88dfe-c000.txt
part-00004-b2128a57-30c4-4114-9ac0-344b81f88dfe-c000.txt
part-00005-b2128a57-30c4-4114-9ac0-344b81f88dfe-c000.txt
part-00006-b2128a57-30c4-4114-9ac0-344b81f88dfe-c000.txt
part-00007-b2128a57-30c4-4114-9ac0-344b81f88dfe-c000.txt
part-00008-b2128a57-30c4-4114-9ac0-344b81f88dfe-c000.txt
part-00009-b2128a57-30c4-4114-9ac0-344b81f88dfe-c000.txt
part-00010-b2128a57-30c4-4114-9ac0-344b81f88dfe-c000.txt
part-00011-b2128a57-30c4-4114-9ac0-344b81f88dfe-c000.txt
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

Edges are stored as two Vertex ID values delimited by tab.
First Vertex identifies host that contains the link.
Second Vertex is host pointed by that link.
Edges are sorted by first Vertex ID using int (not string) comparison logic.
Each file is not continuation of index. Instead, each file contains subset of all Vertices. Therefore we need to get check all files to load all results.


```csv
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


## Create Reversed Edges

Edges store links sorted by site name. Using that sort order we can efficiently (without scanning all files) find all links from particular site to others.

Having all links in one database we can get not only link *from* our site, but also links from other sites *to* our site. Unfortunately, with one set of files we cannot do that without scanning all the files.

We can solve that problem by creating other files where data is stored and sorted in reversed order. In these new edge files we want host be first and second is site that has link pointing to host.


We can use small Python script to reverse Edges fils and save results into `data/edges_reversed` folder:

```bash
cd data && python3 reverse.py 
```

Data is not sorted after Vertices are reversed.
Run script to sort files. **Sorting can take several hours**.

```
$./sort.sh  
```

Now we have about 55 GB of reversed edges in `.txt` files.

```
./data/edges_reversed:
part-00000-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00001-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00002-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00003-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00004-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00005-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00006-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00007-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00008-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00009-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00010-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00011-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00012-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00013-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00014-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00015-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00016-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00017-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00018-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00019-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00020-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00021-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00022-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
part-00023-6318a1c8-2100-4c30-9650-71839e048ef2-c000.txt
```

This is sample data from Sorted Reversed Edges file 

```csv
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

## Create Offset files

Vertices and Edges are stored in sorted files. We can use binary search to find data in file wi want we want with O(log n) time.
But we can do it even faster. We can precompute offsets for every X lines and get data we want with single get operation using these offsets.
It is especially helpful in IO constrain environment as AWS Lambda.

Command line to run offsets precalculation:


```bash
go run cmd/indexer/main.go data
```

You should see 3 files updated:

```
./data:
edges-reversed.offsets.txt
edges.offsets.txt
vertices.offsets.txt
```

## Pushing files to S3

This commandline uploads all files to S3 bucket `common-crawl-hosts` using AWS credentials from profile `commoncrawler`. These credentials need to get write access to bucket.


```bash
aws s3 cp data  s3://common-crawl-hosts/ --recursive --exclude "*" --include "*.txt" --profile commoncrawler
```

