package server

import (
	"math/rand"
)

var seq1 = []Chord{
	{
			ChordSymbol: "Eb",
			MidiValues: []uint8{39,55,58,63,70,82},
	},
	{
			ChordSymbol: "Gb7(addb6)",
			MidiValues: []uint8{30,42,58,62,64,73,78,82},
	},
	{
			ChordSymbol: "C7(b13)",
			MidiValues: []uint8{36,48,58,64,68,72,80},
	},
	{
			ChordSymbol: "Gb13(#11)(9)",
			MidiValues: []uint8{30,42,52,58,63,68,72,80},
	},
	{
			ChordSymbol: "C/Bb",
			MidiValues: []uint8{34,46,52,60,64,67,76},
	},
	{
			ChordSymbol: "A(addb2)/C#",
			MidiValues: []uint8{37,49,57,58,61,64,69,76},
	},
	{
			ChordSymbol: "Dsus4(addb2)",
			MidiValues: []uint8{38,38,50,57,62,63,67,69,74},
	},
	{
			ChordSymbol: "D/F#",
			MidiValues: []uint8{30,54,57,62,66,69,78},
	},
	{
			ChordSymbol: "G-/Bb",
			MidiValues: []uint8{34,50,58,67,74},
	},
	{
			ChordSymbol: "D7/F#",
			MidiValues: []uint8{30,42,57,62,66,69,72,81,84},
	},
	{
			ChordSymbol: "G7sus",
			MidiValues: []uint8{31,43,53,55,60,62,67,74},
	},
	{
			ChordSymbol: "A/E",
			MidiValues: []uint8{40,57,61,69,73,76},
	},
	{
			ChordSymbol: "D7sus",
			MidiValues: []uint8{38,62,67,69,72,74,81},
	},
	{
			ChordSymbol: "D",
			MidiValues: []uint8{38,50,57,66,69,74,78,81},
	},
	{
			ChordSymbol: "G-/Bb",
			MidiValues: []uint8{34,46,62,67,70,79},
	},
	{
			ChordSymbol: "G+(addb2)",
			MidiValues: []uint8{31,43,51,59,63,67,68,71,75,79},
	},
	{
			ChordSymbol: "C-Maj7",
			MidiValues: []uint8{36,48,63,67,71,79},
	},
	{
			ChordSymbol: "C7/G",
			MidiValues: []uint8{31,43,52,58,60,67,76},
	},
	{
			ChordSymbol: "F-",
			MidiValues: []uint8{29,41,53,56,60,68,77},
	},
	{
			ChordSymbol: "A(addb2)/C#",
			MidiValues: []uint8{37,49,61,69,70,73,76},
	},
	{
			ChordSymbol: "Dsus4",
			MidiValues: []uint8{38,50,57,62,67,69,74},
	},
	{
			ChordSymbol: "D7(b13)",
			MidiValues: []uint8{38,50,54,60,62,70},
	},
	{
			ChordSymbol: "D7(#9)",
			MidiValues: []uint8{38,54,60,65,69},
	},
	{
			ChordSymbol: "F#dim/C",
			MidiValues: []uint8{36,48,54,57,60,66,69},
	},
	{
			ChordSymbol: "G7sus",
			MidiValues: []uint8{31,43,53,60,62,67},
	},
	{
			ChordSymbol: "E7(#5)",
			MidiValues: []uint8{40,52,60,62,64,68,76},
	},
	{
			ChordSymbol: "F-",
			MidiValues: []uint8{29,53,60,65,68,72,77},
	},
	{
			ChordSymbol: "CbMaj7(#5)",
			MidiValues: []uint8{35,47,63,67,70,75},
	},
	{
			ChordSymbol: "E-11(9)",
			MidiValues: []uint8{40,54,55,62,66,69,78,81},
	},
	{
			ChordSymbol: "Fb9(#11)",
			MidiValues: []uint8{40,52,62,66,68,70,74,78,82},
	},
	{
			ChordSymbol: "Eb6(9)",
			MidiValues: []uint8{39,46,55,58,60,65,70,77,82},
	},
	{
			ChordSymbol: "Fdim/Ab",
			MidiValues: []uint8{32,44,53,56,59,65,71,80,83},
	},
	{
			ChordSymbol: "G/D",
			MidiValues: []uint8{38,50,59,67,74,79,83},
	},
	{
			ChordSymbol: "Bbdim7",
			MidiValues: []uint8{34,46,55,64,73,82},
	},
	{
			ChordSymbol: "Bb-7(b5)",
			MidiValues: []uint8{34,46,61,64,68,76,80},
	},
	{
			ChordSymbol: "Cb-Maj7",
			MidiValues: []uint8{35,47,54,62,70,74,78},
	},
	{
			ChordSymbol: "Db13(#11)(9)(omit3)",
			MidiValues: []uint8{37,49,59,67,70,75},
	},
	{
			ChordSymbol: "D7(b9)(b5)",
			MidiValues: []uint8{38,54,60,68,72,75},
	},
	{
			ChordSymbol: "F#dim7",
			MidiValues: []uint8{30,42,57,60,63,66,72},
	},
	{
			ChordSymbol: "Gsus4(addb2)",
			MidiValues: []uint8{31,43,50,55,56,60,62,67,74},
	},
	{
			ChordSymbol: "Gb/Bb",
			MidiValues: []uint8{34,46,61,66,70,73},
	},
	{
			ChordSymbol: "G7(#9)(#5)",
			MidiValues: []uint8{31,43,53,59,63,67,70},
	},
	{
			ChordSymbol: "G7(b13)",
			MidiValues: []uint8{31,43,53,55,59,63,71},
	},
	{
			ChordSymbol: "A-11(b5)",
			MidiValues: []uint8{33,45,55,60,62,63,67},
	},
	{
			ChordSymbol: "Ab-6(9)",
			MidiValues: []uint8{32,44,51,53,59,63,65,70,75},
	},
	{
			ChordSymbol: "Db9",
			MidiValues: []uint8{37,53,59,63,68,71},
	},
	{
			ChordSymbol: "G7/D",
			MidiValues: []uint8{38,55,59,62,65,67,71,79},
	},
	{
			ChordSymbol: "G+(addb2)",
			MidiValues: []uint8{31,43,55,59,63,67,68,71,75},
	},
	{
			ChordSymbol: "C-",
			MidiValues: []uint8{36,48,55,63,72,75,79},
	},
	{
			ChordSymbol: "Adim",
			MidiValues: []uint8{33,45,60,63,69},
	},
	{
			ChordSymbol: "D7/C",
			MidiValues: []uint8{36,60,62,66,69,78,81},
	},
	{
			ChordSymbol: "Gsus4(addb2)",
			MidiValues: []uint8{31,55,56,60,62,67,74,79},
	},
	{
			ChordSymbol: "G+(addb2)",
			MidiValues: []uint8{31,43,59,67,68,71,75,79},
	},
	{
			ChordSymbol: "B+",
			MidiValues: []uint8{35,47,51,55,59,63,71,75},
	},
	{
			ChordSymbol: "Ab-/Bb",
			MidiValues: []uint8{34,46,56,59,63,68,75,80},
	},
	{
			ChordSymbol: "G7(b13)",
			MidiValues: []uint8{31,43,53,55,59,63,67,75,79},
	},
	{
			ChordSymbol: "C-Maj7",
			MidiValues: []uint8{36,51,59,67,71,75},
	},
	{
			ChordSymbol: "C-7",
			MidiValues: []uint8{36,51,55,58,67,75},
	},
	{
			ChordSymbol: "G7(#9)(#5)",
			MidiValues: []uint8{31,43,53,55,59,63,70},
	},
	{
			ChordSymbol: "Csus4",
			MidiValues: []uint8{36,53,60,67,72},
	},
	{
			ChordSymbol: "C+",
			MidiValues: []uint8{36,52,56,60,64,68},
	},
	{
			ChordSymbol: "F-6",
			MidiValues: []uint8{29,41,50,53,56,60,68},
	},
	{
			ChordSymbol: "Asus4",
			MidiValues: []uint8{33,45,57,62,64,69},
	},
	{
			ChordSymbol: "AbMaj7(#9)",
			MidiValues: []uint8{32,44,51,55,59,60,63,67},
	},
	{
			ChordSymbol: "C-(addb6)/G",
			MidiValues: []uint8{31,55,56,60,67},
	},
	{
			ChordSymbol: "B+",
			MidiValues: []uint8{35,47,55,59,63,67},
	},
	{
			ChordSymbol: "C-(addb6)",
			MidiValues: []uint8{36,51,55,56,60,63},
	},
	{
			ChordSymbol: "F(add2)/A",
			MidiValues: []uint8{33,45,60,65,67,69},
	},
	{
			ChordSymbol: "D+",
			MidiValues: []uint8{38,54,58,62,66},
	},
	{
			ChordSymbol: "G-9",
			MidiValues: []uint8{31,43,57,58,62,65},
	},
	{
			ChordSymbol: "G-(add4)",
			MidiValues: []uint8{31,43,50,55,58,60,62},
	},
	{
			ChordSymbol: "AbMaj7(#9)",
			MidiValues: []uint8{32,44,51,55,59,60,63,67},
	},
	{
			ChordSymbol: "Fdim/Ab",
			MidiValues: []uint8{32,44,53,56,59,65},
	},
	{
			ChordSymbol: "Eb/G",
			MidiValues: []uint8{31,43,51,55,58,63},
	},
	{
			ChordSymbol: "G(addb2)/D",
			MidiValues: []uint8{38,50,55,56,59,62},
	},
	{
			ChordSymbol: "Bdim/F",
			MidiValues: []uint8{29,41,50,53,59,62},
	},
	{
			ChordSymbol: "C-6(9)",
			MidiValues: []uint8{36,51,55,57,62},
	},
	{
			ChordSymbol: "G7(#9)",
			MidiValues: []uint8{31,43,59,62,65,70},
	},
	{
			ChordSymbol: "A-9(b5)",
			MidiValues: []uint8{33,45,60,63,67,71},
	},
	{
			ChordSymbol: "A-7(b5)",
			MidiValues: []uint8{33,45,55,60,63,67},
	},
	{
			ChordSymbol: "F#dim/A",
			MidiValues: []uint8{33,45,54,60,66,69},
	},
	{
			ChordSymbol: "G-/Bb",
			MidiValues: []uint8{34,46,50,55,58,62,70,74},
	},
	{
			ChordSymbol: "D(addb2)/A",
			MidiValues: []uint8{33,45,54,62,63,66,69},
	},
	{
			ChordSymbol: "D+(addb2)",
			MidiValues: []uint8{38,50,54,58,62,63,66,70},
	},
	{
			ChordSymbol: "F#dim7",
			MidiValues: []uint8{30,54,60,63,66,69,72},
	},
	{
			ChordSymbol: "Gsus4",
			MidiValues: []uint8{31,43,55,60,62,67,74},
	},
	{
			ChordSymbol: "F-13(9)(b5)",
			MidiValues: []uint8{29,41,51,55,56,59,62,67,74},
	},
	{
			ChordSymbol: "G(addb2)/D",
			MidiValues: []uint8{38,55,56,59,62,67,71,74},
	},
	{
			ChordSymbol: "B+(addb2)",
			MidiValues: []uint8{35,47,55,63,67,71,72,75,79},
	},
	{
			ChordSymbol: "C-Maj9",
			MidiValues: []uint8{36,51,59,62,67,71,79},
	},
	{
			ChordSymbol: "Bbdim/Fb",
			MidiValues: []uint8{28,40,52,58,61,64,70,73,82},
	},
	{
			ChordSymbol: "C(addb2)/E",
			MidiValues: []uint8{40,55,60,61,64,67,72,79},
	},
	{
			ChordSymbol: "F-",
			MidiValues: []uint8{29,53,56,60,68,77,80},
	},
	{
			ChordSymbol: "Bb9sus",
			MidiValues: []uint8{34,46,56,60,63,65,68,72,77,80},
	},
	{
			ChordSymbol: "G7/B",
			MidiValues: []uint8{35,47,55,62,65,74,83},
	},
	{
			ChordSymbol: "Db+",
			MidiValues: []uint8{37,49,57,65,69,73,81},
	},
	{
			ChordSymbol: "Gb-(add2)/A",
			MidiValues: []uint8{33,45,54,61,66,68,69,73,78},
	},
	{
			ChordSymbol: "A-6",
			MidiValues: []uint8{33,57,60,64,66,72,76,81},
	},
	{
			ChordSymbol: "F+",
			MidiValues: []uint8{29,53,61,69,73,77},
	},
	{
			ChordSymbol: "Dsus4",
			MidiValues: []uint8{38,50,67,69,74},
	},
	{
			ChordSymbol: "D+",
			MidiValues: []uint8{38,50,54,58,62,66,74},
	},
	{
			ChordSymbol: "F#dim/A",
			MidiValues: []uint8{33,45,57,60,66,72},
	},
	{
			ChordSymbol: "G-/D",
			MidiValues: []uint8{38,55,58,62,70},
	},
	{
			ChordSymbol: "Bb-(add2)",
			MidiValues: []uint8{34,46,53,58,60,61,65,70},
	},
	{
			ChordSymbol: "G(addb2)/D",
			MidiValues: []uint8{38,50,55,56,59,62},
	},
	{
			ChordSymbol: "G+",
			MidiValues: []uint8{31,43,51,55,59,63},
	},
	{
			ChordSymbol: "E+",
			MidiValues: []uint8{40,56,60,64},
	},
	{
			ChordSymbol: "F-6(9)",
			MidiValues: []uint8{29,41,48,56,60,62,67,72},
	},
	{
			ChordSymbol: "Eb/Bb",
			MidiValues: []uint8{34,51,55,63,67,70},
	},
	{
			ChordSymbol: "Bb-13(9)(b5)",
			MidiValues: []uint8{34,46,52,56,61,64,72,79},
	},
	{
			ChordSymbol: "Edim7",
			MidiValues: []uint8{40,52,58,61,64,67,76},
	},
	{
			ChordSymbol: "F-/C",
			MidiValues: []uint8{36,48,56,60,65,68,72,77},
	},
	{
			ChordSymbol: "G7sus",
			MidiValues: []uint8{31,55,65,67,72,74,79},
	},
	{
			ChordSymbol: "Bdim/D",
			MidiValues: []uint8{38,50,59,65,74,83},
	},
	{
			ChordSymbol: "C-(b6)(9)",
			MidiValues: []uint8{36,43,51,55,56,62,67,74,79},
	},
	{
			ChordSymbol: "C7/E",
			MidiValues: []uint8{28,40,55,58,60,67,72,76,79},
	},
	{
			ChordSymbol: "Edim",
			MidiValues: []uint8{40,52,58,64,67,70,79},
	},
	{
			ChordSymbol: "F-6(9)",
			MidiValues: []uint8{29,41,56,60,62,67,72,79,84},
	},
	{
			ChordSymbol: "CbMaj7(#9)",
			MidiValues: []uint8{35,51,54,58,66,70,74,82},
	},
	{
			ChordSymbol: "Db13(b9)",
			MidiValues: []uint8{37,49,58,59,62,65,70,77,82},
	},
	{
			ChordSymbol: "Eb-(addb6)",
			MidiValues: []uint8{39,54,58,59,63,66,70,78,82},
	},
	{
			ChordSymbol: "D7sus",
			MidiValues: []uint8{38,62,67,69,72,74,81},
	},
	{
			ChordSymbol: "Bdim",
			MidiValues: []uint8{35,47,53,62,65,71,74},
	},
	{
			ChordSymbol: "C-(b6)(9)",
			MidiValues: []uint8{36,43,51,55,56,62,67,74,79},
	},
	{
			ChordSymbol: "Gb13(b9)",
			MidiValues: []uint8{30,42,52,55,58,63,70,75},
	},
	{
			ChordSymbol: "F-6(9)",
			MidiValues: []uint8{29,53,60,62,67,68,72,79},
	},
	{
			ChordSymbol: "Bb7/F",
			MidiValues: []uint8{29,53,62,65,68,70,74,77},
	},
}
var seq2 = []Chord{
	{
			ChordSymbol: "Eb",
			MidiValues: []uint8{39,46,51,55,58,67},
	},
	{
			ChordSymbol: "CbMaj7(#5)",
			MidiValues: []uint8{35,47,63,67,70},
	},
	{
			ChordSymbol: "Bbsus4(add2)",
			MidiValues: []uint8{34,58,60,63,65,70},
	},
	{
			ChordSymbol: "C-7(addb6)",
			MidiValues: []uint8{36,48,55,56,58,63,70},
	},
	{
			ChordSymbol: "Gdim/Db",
			MidiValues: []uint8{37,49,55,58,61,67,70,73},
	},
	{
			ChordSymbol: "Ab(add2)",
			MidiValues: []uint8{32,44,51,60,63,68,70,72},
	},
	{
			ChordSymbol: "C-",
			MidiValues: []uint8{36,43,51,60,67,75},
	},
	{
			ChordSymbol: "Ebdim9",
			MidiValues: []uint8{39,51,54,57,60,65},
	},
	{
			ChordSymbol: "G-(add2)",
			MidiValues: []uint8{31,55,62,67,69,70,74},
	},
	{
			ChordSymbol: "Ab-13(11)(9)",
			MidiValues: []uint8{32,44,54,58,59,65,73},
	},
	{
			ChordSymbol: "C-Maj7",
			MidiValues: []uint8{36,48,59,63,67,75,83},
	},
	{
			ChordSymbol: "Bbdim11",
			MidiValues: []uint8{34,46,61,64,67,70,75},
	},
	{
			ChordSymbol: "F-6(9)",
			MidiValues: []uint8{29,53,60,62,68,72,79},
	},
	{
			ChordSymbol: "Cbsus4(add2)",
			MidiValues: []uint8{35,47,64,66,71,73,78},
	},
	{
			ChordSymbol: "F+(addb2)",
			MidiValues: []uint8{29,41,57,61,65,66,69,73,77},
	},
	{
			ChordSymbol: "Bb-",
			MidiValues: []uint8{34,46,61,65,70,73},
	},
	{
			ChordSymbol: "Db-(add2)",
			MidiValues: []uint8{37,49,64,68,73,75,76},
	},
	{
			ChordSymbol: "Fb-",
			MidiValues: []uint8{28,52,55,59,67,76},
	},
	{
			ChordSymbol: "E+",
			MidiValues: []uint8{28,40,56,60,64,68,72},
	},
	{
			ChordSymbol: "F-6",
			MidiValues: []uint8{29,41,50,56,65,68,72,77},
	},
	{
			ChordSymbol: "FbMaj9",
			MidiValues: []uint8{40,56,59,63,66,71,75},
	},
	{
			ChordSymbol: "C+",
			MidiValues: []uint8{36,52,56,60,64,68,72},
	},
	{
			ChordSymbol: "Adim/C",
			MidiValues: []uint8{36,51,57,60,63,69,72},
	},
	{
			ChordSymbol: "Ab-",
			MidiValues: []uint8{32,44,56,59,63,71},
	},
	{
			ChordSymbol: "Bb13sus(9)",
			MidiValues: []uint8{34,46,55,56,60,63,70},
	},
	{
			ChordSymbol: "C-/G",
			MidiValues: []uint8{31,43,60,63,67},
	},
	{
			ChordSymbol: "A(addb6)",
			MidiValues: []uint8{33,45,52,61,64,65,69,76},
	},
	{
			ChordSymbol: "Ab(add2)/C",
			MidiValues: []uint8{36,51,56,58,60,63,72,75},
	},
	{
			ChordSymbol: "D7(b13)(b9)",
			MidiValues: []uint8{38,50,60,63,66,70},
	},
	{
			ChordSymbol: "D/F#",
			MidiValues: []uint8{30,42,54,57,62,69},
	},
	{
			ChordSymbol: "F#dim7(b13)",
			MidiValues: []uint8{30,54,60,63,66,69,74},
	},
	{
			ChordSymbol: "G7sus",
			MidiValues: []uint8{31,43,60,62,65,67,74},
	},
	{
			ChordSymbol: "Fdim/Ab",
			MidiValues: []uint8{32,44,59,65,68,71},
	},
	{
			ChordSymbol: "G+(addb2)",
			MidiValues: []uint8{31,43,51,55,59,63,67,68,71,75},
	},
	{
			ChordSymbol: "E7(#9)",
			MidiValues: []uint8{28,40,56,59,62,67,74},
	},
	{
			ChordSymbol: "FbMaj7(#5)",
			MidiValues: []uint8{28,40,56,60,63},
	},
	{
			ChordSymbol: "Adim/Eb",
			MidiValues: []uint8{39,51,57,60,69},
	},
	{
			ChordSymbol: "Ab-/Cb",
			MidiValues: []uint8{35,51,56,59,63,68},
	},
	{
			ChordSymbol: "Fb-/Cb",
			MidiValues: []uint8{35,47,55,64,71,79},
	},
	{
			ChordSymbol: "Db+(addb2)",
			MidiValues: []uint8{37,49,53,57,61,62,65,69},
	},
	{
			ChordSymbol: "Gb(b6)(9)",
			MidiValues: []uint8{30,42,49,58,61,62,68,73},
	},
	{
			ChordSymbol: "F+",
			MidiValues: []uint8{29,41,57,61,65,73},
	},
	{
			ChordSymbol: "Dsus4",
			MidiValues: []uint8{38,62,67,69,74},
	},
	{
			ChordSymbol: "Ab9(#11)(omit3)",
			MidiValues: []uint8{32,44,54,62,66,70,74,78},
	},
	{
			ChordSymbol: "Eb(add2)",
			MidiValues: []uint8{39,46,55,58,63,65,67,75,79},
	},
	{
			ChordSymbol: "C7",
			MidiValues: []uint8{36,48,58,60,67,76,79},
	},
	{
			ChordSymbol: "A/C#",
			MidiValues: []uint8{37,49,52,61,69,76},
	},
	{
			ChordSymbol: "Ab/C",
			MidiValues: []uint8{36,48,51,56,60,63,72,75},
	},
	{
			ChordSymbol: "Eb-/Bb",
			MidiValues: []uint8{34,51,54,58,66,75},
	},
	{
			ChordSymbol: "Eb13sus(b9)(add3)",
			MidiValues: []uint8{39,56,61,64,67,72},
	},
	{
			ChordSymbol: "C7",
			MidiValues: []uint8{36,52,58,60,67,76},
	},
	{
			ChordSymbol: "Adim",
			MidiValues: []uint8{33,45,60,63,69,75},
	},
	{
			ChordSymbol: "D/A",
			MidiValues: []uint8{33,45,54,57,62,66},
	},
	{
			ChordSymbol: "G-/Bb",
			MidiValues: []uint8{34,50,55,58,67},
	},
	{
			ChordSymbol: "Cdim/Gb",
			MidiValues: []uint8{30,42,51,60,66},
	},
	{
			ChordSymbol: "G7(b9)(b5)",
			MidiValues: []uint8{31,43,53,56,59,61,65},
	},
	{
			ChordSymbol: "G7/B",
			MidiValues: []uint8{35,47,53,55,59,62,67},
	},
	{
			ChordSymbol: "C-6(9)",
			MidiValues: []uint8{36,51,55,57,62,67,74},
	},
	{
			ChordSymbol: "AbMaj9",
			MidiValues: []uint8{32,44,51,55,58,60,63,67,75},
	},
	{
			ChordSymbol: "Edim7(b13)",
			MidiValues: []uint8{40,52,58,61,64,67,72},
	},
	{
			ChordSymbol: "A7sus",
			MidiValues: []uint8{33,45,62,64,67,69,76},
	},
	{
			ChordSymbol: "Ab(add2)/C",
			MidiValues: []uint8{36,51,56,58,63,72,75},
	},
	{
			ChordSymbol: "C7(#9)(b5)",
			MidiValues: []uint8{36,52,58,63,66,70,75,78},
	},
	{
			ChordSymbol: "Bbdim9",
			MidiValues: []uint8{34,58,61,64,67,72,79},
	},
	{
			ChordSymbol: "F-(add2)/C",
			MidiValues: []uint8{36,53,55,56,60,68,77,80},
	},
	{
			ChordSymbol: "E+",
			MidiValues: []uint8{28,40,52,56,60,64,68,72,80,84},
	},
	{
			ChordSymbol: "Db7(#9)(#5)",
			MidiValues: []uint8{37,53,59,64,69,73,81},
	},
	{
			ChordSymbol: "Gb7sus",
			MidiValues: []uint8{30,42,52,59,61,66,73,78},
	},
	{
			ChordSymbol: "Eb-6(9)",
			MidiValues: []uint8{39,53,54,58,60,65,70,77},
	},
	{
			ChordSymbol: "C7",
			MidiValues: []uint8{36,52,58,60,67,76},
	},
	{
			ChordSymbol: "A7/C#",
			MidiValues: []uint8{37,52,55,57,61,64,73},
	},
	{
			ChordSymbol: "Dsus4",
			MidiValues: []uint8{38,50,55,57,62,69,74},
	},
	{
			ChordSymbol: "D(addb2)/F#",
			MidiValues: []uint8{30,54,62,63,66,69},
	},
	{
			ChordSymbol: "G-(add2)/Bb",
			MidiValues: []uint8{34,46,55,57,58,62,67,70,74,79},
	},
	{
			ChordSymbol: "F#dim/C",
			MidiValues: []uint8{36,48,54,57,60,66,69,72,78},
	},
	{
			ChordSymbol: "G-(b6)(9)",
			MidiValues: []uint8{31,43,50,58,62,63,69,74},
	},
	{
			ChordSymbol: "F-/Ab",
			MidiValues: []uint8{32,44,60,65,68,77,84},
	},
	{
			ChordSymbol: "D(addb2)/F#",
			MidiValues: []uint8{30,42,57,66,69,74,75,78,81},
	},
	{
			ChordSymbol: "F#dim7",
			MidiValues: []uint8{30,42,51,57,60,66,72,84},
	},
	{
			ChordSymbol: "G-(b6)(9)",
			MidiValues: []uint8{31,43,50,58,62,63,69,74,81},
	},
	{
			ChordSymbol: "Ab(add2)/C",
			MidiValues: []uint8{36,51,56,58,63,68,72,80,84},
	},
	{
			ChordSymbol: "C/Bb",
			MidiValues: []uint8{34,46,52,55,60,64,76},
	},
	{
			ChordSymbol: "Edim/Bb",
			MidiValues: []uint8{34,46,55,58,64,67,76,79},
	},
	{
			ChordSymbol: "B+(addb2)",
			MidiValues: []uint8{35,59,60,63,67,71,75,79},
	},
	{
			ChordSymbol: "Fb-Maj7",
			MidiValues: []uint8{40,52,55,63,67,71,75},
	},
	{
			ChordSymbol: "Gb13(#11)(9)",
			MidiValues: []uint8{30,42,52,56,58,63,72,75,80},
	},
	{
			ChordSymbol: "Edim",
			MidiValues: []uint8{40,52,58,67,76},
	},
	{
			ChordSymbol: "F-6(9)",
			MidiValues: []uint8{29,53,60,62,67,68,72,79},
	},
	{
			ChordSymbol: "C7sus(b9)",
			MidiValues: []uint8{36,48,58,61,65,67,72},
	},
	{
			ChordSymbol: "C7/E",
			MidiValues: []uint8{40,52,58,60,64,67,72},
	},
	{
			ChordSymbol: "Edim7",
			MidiValues: []uint8{28,52,58,61,64,67},
	},
	{
			ChordSymbol: "Edim7(b13)",
			MidiValues: []uint8{40,52,58,61,64,67,72},
	},
	{
			ChordSymbol: "F-6",
			MidiValues: []uint8{29,41,50,56,60,65,68,72},
	},
	{
			ChordSymbol: "C7/G",
			MidiValues: []uint8{31,55,60,64,70},
	},
	{
			ChordSymbol: "F-6(9)",
			MidiValues: []uint8{29,41,55,56,60,62,67},
	},
	{
			ChordSymbol: "Abdim/D",
			MidiValues: []uint8{38,50,56,59,62,68},
	},
	{
			ChordSymbol: "Eb(add2)",
			MidiValues: []uint8{39,46,51,53,55,58,63},
	},
	{
			ChordSymbol: "C-Maj7",
			MidiValues: []uint8{36,43,51,59,67},
	},
	{
			ChordSymbol: "G7(#9)(#5)",
			MidiValues: []uint8{31,43,53,58,59,63,67},
	},
	{
			ChordSymbol: "A-11(9)(b5)",
			MidiValues: []uint8{33,45,55,60,62,63,67,71},
	},
	{
			ChordSymbol: "F-9(b5)",
			MidiValues: []uint8{29,41,51,55,56,59,63,67},
	},
	{
			ChordSymbol: "F-7(b5)",
			MidiValues: []uint8{29,41,51,56,59,63,68},
	},
	{
			ChordSymbol: "Eb/Bb",
			MidiValues: []uint8{34,46,55,63,67,70,75},
	},
	{
			ChordSymbol: "C/E",
			MidiValues: []uint8{28,52,55,60,67,76},
	},
	{
			ChordSymbol: "Edim7",
			MidiValues: []uint8{28,52,61,67,70,76},
	},
	{
			ChordSymbol: "F-6(9)",
			MidiValues: []uint8{29,41,56,60,62,67,72,79},
	},
	{
			ChordSymbol: "F-13(11)",
			MidiValues: []uint8{29,41,51,56,62,65,70,74,77},
	},
	{
			ChordSymbol: "Ab-(add2)/Eb",
			MidiValues: []uint8{39,56,58,59,63,68,71,80,83},
	},
	{
			ChordSymbol: "F13(#11)(9)(omit3)",
			MidiValues: []uint8{29,41,51,55,62,71,79,83},
	},
	{
			ChordSymbol: "Bb-11(9)",
			MidiValues: []uint8{34,46,61,68,72,75,84},
	},
	{
			ChordSymbol: "Db-Maj9",
			MidiValues: []uint8{37,49,52,60,63,68,72,75,84},
	},
	{
			ChordSymbol: "G7(#9)",
			MidiValues: []uint8{31,43,59,62,65,70,74,82},
	},
	{
			ChordSymbol: "E7(#9)(b5)",
			MidiValues: []uint8{40,56,58,62,67,70,79},
	},
	{
			ChordSymbol: "Db7(#5)",
			MidiValues: []uint8{37,49,57,59,61,65,69,77,81},
	},
	{
			ChordSymbol: "Gb+",
			MidiValues: []uint8{30,42,54,58,62,66,70,78},
	},
	{
			ChordSymbol: "Fb9(#11)",
			MidiValues: []uint8{28,40,56,62,66,70,74,82},
	},
	{
			ChordSymbol: "Eb-6(9)",
			MidiValues: []uint8{39,46,53,54,58,60,65,70,77,82},
	},
	{
			ChordSymbol: "D7(b9)",
			MidiValues: []uint8{38,50,54,63,72,78},
	},
	{
			ChordSymbol: "B+(addb2)",
			MidiValues: []uint8{35,51,55,63,71,72,75,79},
	},
	{
			ChordSymbol: "Bbsus4(addb2)",
			MidiValues: []uint8{34,46,58,59,63,65,70,77},
	},
	{
			ChordSymbol: "A-7(b5)",
			MidiValues: []uint8{33,45,60,63,67,72,75},
	},
	{
			ChordSymbol: "F-11(9)(b5)",
			MidiValues: []uint8{29,41,56,59,63,67,70,75},
	},
	{
			ChordSymbol: "Bb+(addb2)",
			MidiValues: []uint8{34,50,54,58,59,62,66,70},
	},
	{
			ChordSymbol: "Cdim7",
			MidiValues: []uint8{36,48,54,57,63,66,72},
	},
	{
			ChordSymbol: "F#dim/A",
			MidiValues: []uint8{33,45,54,60,69},
	},
	{
			ChordSymbol: "G-(b6)(9)",
			MidiValues: []uint8{31,43,58,62,63,69,74},
	},
	{
			ChordSymbol: "Ab(add2)/C",
			MidiValues: []uint8{36,48,63,68,70,72},
	},
}

var dummySeq1 = func() []Chord {
	dummySeq := make([]Chord, len(seq1))
	for i, chord := range seq1 {
		values := make([]uint8, len(chord.MidiValues))
		copy(values, chord.MidiValues)
		dummySeq[i] = Chord{chord.ChordSymbol, values}
	}
	return dummySeq
}()

var dummySeq2 = func() []Chord {
	dummySeq := make([]Chord, len(seq2))
	for i, chord := range seq1 {
		values := make([]uint8, len(chord.MidiValues))
		copy(values, chord.MidiValues)
		dummySeq[i] = Chord{chord.ChordSymbol, values}
	}
	return dummySeq
}()

func getDummyData() []Chord {
	var randomSeq []Chord
	if rand.Float64() < 0.5 {
		randomSeq = dummySeq1
	} else {
		randomSeq = dummySeq2
	}
	startIndex := rand.Intn(len(randomSeq)-24+1)
	slicedSeq := randomSeq[startIndex : startIndex+24]
	return append([]Chord{}, slicedSeq...)
}
