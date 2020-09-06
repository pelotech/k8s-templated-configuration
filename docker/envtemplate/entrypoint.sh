#! /bin/sh
for path in $@; do
  for file in $(find $path -type f); do
    echo $file
    envtemplate <$file | sponge $file
  done
done
