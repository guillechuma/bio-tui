# Bio-TUI

**Terminal-native bioinformatics file viewer** — rich, dynamic, and visually pleasing.
Say goodbye to GUIs and browsers. Say hello to the command line, where bioinformaticians live.

## Overview

**Bio-TUI** is an open-source, cross-platform TUI (Text User Interface) for exploring the most common bioinformatics file formats — right in your terminal.

## Supported Formats

- **BAM / CRAM** — coverage, pileup, min/max/mean stats
- **VCF / BCF** — variant table, filter highlights
- **FASTA / FASTQ** — sequence chunks, GC%, quality visualization
- **GFF3 / GTF** — gene annotations, exon/CDS tracks

## Features

- **Instant jumps** to genomic regions, genes, or read IDs
- **Coverage tracks** with Unicode graphics
- **Pileup view** for deep inspection
- **Variant table** with filters and impact coloring
- **Annotation lanes** for GFF/GTF data
- **Export to PNG/JSON** for reports or sharing
- Works **entirely offline** — single static binary

## Installation

### From source

```bash
go install github.com/guillechuma/bio-tui/cmd/bio-tui@latest
```

### Prebuilt binaries

Coming soon

## Getting Started

Coming soon

## License

Bio-TUI is released under the [MIT License](https://opensource.org/licenses/MIT).
