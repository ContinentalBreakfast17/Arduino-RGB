#!/bin/bash
# Run this to generate "full.html" for in browser testing.
# However, this is currently not functional as there is pre-loaded data in this mug 

echo "<style>" > full.html
cat styles.css >> full.html
echo "" >> full.html
echo "</style>" >> full.html

cat web.html >> full.html
echo "" >> full.html

echo "<script src=\"https://cdn.jsdelivr.net/npm/vue\"></script>" >> full.html

echo "<script>" >> full.html
cat vue/app.js >> full.html
echo "" >> full.html
echo "</script>" >> full.html