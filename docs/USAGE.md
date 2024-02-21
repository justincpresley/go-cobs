# Usage

## Configuration

The following are variables of the `Config` structure. All variables are applicable to all types. If you are frightened by the options, feel free to use the suggested defaults or create an empty structure (which acts like a very simple implementation).

* **Type**:
  * Description: Type of COBS that is used for all API functions.
  * Possible Values: Any `EncoderType` value.
  * Suggested Default: `Native`.
* **SpecialByte**:
  * Description: The particular byte value to be "encoded away".
  * Possible Values: Any byte value such as `0x00`.
  * Suggested Default: No default.
* **Delimiter**:
  * Description: Whether to include a delimiter: the special byte value placed at the end of the encoded slice to mark the end.
  * Possible Values: `true` and `false`.
  * Suggested Default: No default.
* **EndingSave**:
  * Description: Whether to save a byte when encountering "max" flag(s) occurring at the end of the slice.
  * Possible Values: `true` and `false`.
  * Suggested Default: `true`.

## Types

The following are the types (extensions) that are included. You can think of types as different versions COBS. Each has their own advantages, disadvantages, and use case.

Individual Types:

* [X] **Native** `(COBS)`:
  * Description: This is the natural and default algorithm that was first discovered.
  * Notes: Relatively stable overhead. The easiest to implement.
* [X] **Reduced** `(COBS/R)`:
  * Description: Saves a byte by replacing the last flag with the last character if is appropriate!
  * Notes: Potentially saves a byte of overhead. Encoding possibly generates no overhead. A massive reduction in coverage from flag-based verification.
* [X] **PairElimination** `(COBS/PE)`:
  * Description: Incorporate flags to represent a "pair" of special bytes.
  * Notes: Encoded size can be ~half of the original size. Good for embedded systems. An increase in theoretical overhead, but unlikely to occur.
* [ ] **RunElimination** `(COBS/RE)`:
  * Description: Incorporate flags to represent a "run" of special bytes.
  * Notes: Takes `PairElimination` further.

Combined Types:

* [ ] **PairAndRun** `(COBS/PAR)`:
  * Description: Combines `PairElimination` and `RunElimination`.

Types can be done in **Reversed**:
  * Description: Places the flag at the end of the chunk rather than before, effectively reversing the process.
  * Notes: Allows encoding with no lookahead (more performant most of the time) yet enforces decoding to be done in reverse killing the ability to stream decode.
  * Types:
    * [X] **Reversed** `(RCOBS)`
    * [X] **PairInReverse** `(COBS/PIR)`
    * [ ] **RunInReverse** `(COBS/RIR)`
