# jobs-analyzer

This small tool helps to get some information about specified job roles on US market.

It goes through the most popular job sites in US and aggregates the number of ads for specified role. Currently the only available chart is "Top ten states for specific role".

This may help you to figure out what states to focus on during your job search.

### Example
```
./jobs-analyzer -p "system engineer"
```

### Output
Top 10 states for 'system engineer' role:
* California/CA  —  9495
* Virginia/VA  —  3230
* New York/NY  —  3157
* Texas/TX  —  3065
* Washington/WA  —  2759
* Massachusetts/MA  —  2185
* Illinois/IL  —  1750
* Maryland/MD  —  1725
* Georgia/GA  —  1622
* New Jersey/NJ  —  1570
