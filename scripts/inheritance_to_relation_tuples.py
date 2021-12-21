import json
import argparse
import os
from typing import Dict

rt = {"namespace": "groups", "object": "", "relation": "member", "subject_set": {"namespace": "groups", "object": "", "relation": "member"}}

if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Convert Inheritance to Google Zanzibar Tuples"
    )
    parser.add_argument(
        "inheritance",
        help="an integer for the accumulator",
    )
    parser.add_argument(
        "outfolder",
        help="sum the integers (default: find the max)",
    )
    args = parser.parse_args()
    inheritance: Dict[str, str] = json.load(open(args.inheritance))
    for key in inheritance:
        for role in inheritance[key]:
            rt["object"] = role
            rt["subject_set"]["object"] = key
            with open(
                os.path.join(args.outfolder, "{}_{}.json".format(key, role)), "w"
            ) as f:
                json.dump(rt, f, indent=2)
